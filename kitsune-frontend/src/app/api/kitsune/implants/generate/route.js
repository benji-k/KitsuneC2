import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function POST(req){
    /*const session = await getServerSession(authOptions)
    if (!session) {
        return Response.json({ "error": "Unauthorized" }, { status: 401 })
    }*/

    try{
        const body = await req.json()

        const result = await fetch(API_URL + "/implants/generate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body)
        })

        if (result.status === 200) {
            const blob = await result.blob()
            return new Response(blob)
        } else {
            const error = await result.json()
            return Response.json({ "error": error }, { status: 500 })
        }

    } catch(e){
        return Response.json({ "error": e.message}, { status: 500 })
    }
}