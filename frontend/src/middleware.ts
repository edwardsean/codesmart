import { NextRequest, NextResponse } from "next/server";


export function middleware(req: NextRequest) {
    const access_token = req.cookies.get("access_token")?.value;
    const protected_routes = ["/dashboard"]

    const url = req.nextUrl.clone()

    //if authenticated and visiting again
    console.log("access token: ", access_token)
    console.log("url pathname: ", url.pathname)
    if(access_token && url.pathname == "/") {
        url.pathname = "/dashboard"
        return NextResponse.redirect(url)
    }
    
    //if not authenticated and want to a protected page
    if (!access_token && protected_routes.some((route) => req.nextUrl.pathname.startsWith(route))) {
        const url = new URL("/login", req.url)
        return NextResponse.redirect(url)
    }

    return NextResponse.next()
}

export const config = {
  matcher: ["/dashboard/:path*", "/settings/:path*", "/profile/:path*", "/"],
};