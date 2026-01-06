import { NextRequest, NextResponse } from "next/server";


export async function middleware(req: NextRequest) {
    const refresh_token = req.cookies.get("refresh_token")?.value;
    const public_path = ["/login", "/signup", "/"]


    const url = req.nextUrl.clone()
    const isPublicPath = public_path.includes(url.pathname)

    //if authenticated and visiting again
    if(refresh_token && isPublicPath) {
        url.pathname = "/dashboard"
        return NextResponse.redirect(url)
    }
    
    //if not authenticated and want to a protected page
    if (!refresh_token && !isPublicPath) {
        url.pathname = "/login"
        return NextResponse.redirect(url)
    }

    return NextResponse.next()
}

export const config = {
  matcher: ["/dashboard/:path*", "/settings/:path*", "/profile/:path*", "/"],
};