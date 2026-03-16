import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";


const BACKEND_URL = process.env.BACKEND_URL;

//browser calls next js
//next js server calls api gateway
//api gateway receives refreshtoken header if available, body, headers
//gateway authenticates its access token from header, and then forwards the remaining headers
//gateway forwards to services
//response comes back to next js, and next js forwards to browser
async function handler(req: NextRequest, {params}: {params: Promise<{path: string[]}>}) {
    const { path } = await params;
    const pathString = path.join('/');
    console.log("path string: ", pathString);
    const search = req.nextUrl.search; //preserve query strings like ?code=xxx
    const url = `${BACKEND_URL}/api/v1/${pathString}${search}`;
    console.log("full url: ", url);

    //get refresh token
    const cookieStore = cookies();
    const refreshToken = (await cookieStore).get('refresh_token');

    //forward headers
    const headers : HeadersInit = {
        'Content-Type': req.headers.get('Content-Type') || 'application/json',
    };

    const auth = req.headers.get('Authorization');
    if(auth) headers['Authorization'] = auth;

    //forward refresh token for /auth/refresh and /auth/logout
    if(refreshToken) headers['Cookie'] = `refresh_token=${refreshToken.value}`;

    const body = req.method !== 'GET' && req.method !== 'HEAD' ? await req.text() : undefined;

    const response = await fetch(url, {
        method: req.method,
        headers,
        body,
        redirect: 'manual'
    });

    //for redirects
    if (response.status >= 300 && response.status < 400) {
        const location = response.headers.get('location');
        if (location) {
            return NextResponse.redirect(location);
        }
    }

    const data = await response.text();
    const responseHeaders = new Headers();

    //set the cookie from gateway
    const setCookie = response.headers.get('set-cookie'); //from the response it as Set-Cookie: refresh_token=xxx
    if(setCookie) responseHeaders.set('set-cookie', setCookie);

    responseHeaders.set('Content-Type', response.headers.get('Content-Type') || 'application/json');

    return new NextResponse(data, {
        status: response.status,
        headers: responseHeaders,
    });

}


export const GET     = handler;
export const POST    = handler;
export const PUT     = handler;
export const PATCH   = handler;
export const DELETE  = handler;