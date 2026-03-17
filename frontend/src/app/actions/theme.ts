"use server"
import { cookies } from "next/headers"

export async function setTheme(theme: "light" | "dark") {
    (await cookies()).set("theme", theme, {
        httpOnly: false, //needs to be readable by JS for the toggle
        maxAge: 60 * 60 * 24 * 365, //1 year
        path: "/",
    })
}