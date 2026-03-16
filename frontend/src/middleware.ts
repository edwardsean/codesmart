import { NextRequest, NextResponse } from "next/server";

const public_path = ["/auth/login", "/auth/signup", "/"]

export async function middleware(req: NextRequest) {
    const refresh_token = req.cookies.get("refresh_token")?.value;


    const url = req.nextUrl.clone()
    const isPublicPath = public_path.includes(url.pathname)

    //if authenticated and visiting again
    if(refresh_token && isPublicPath) {
        url.pathname = "/dashboard"
        return NextResponse.redirect(url)
    }
    
    //if not authenticated and want to a protected page
    if (!refresh_token && !isPublicPath) {
        url.pathname = "/auth/login"
        return NextResponse.redirect(url)
    }

    return NextResponse.next()
}

export const config = {
  matcher: ["/dashboard/:path*", "/settings/:path*", "/profile/:path*", "/"],
};