import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function POST(req) {
    //const session = await getServerSession(authOptions)
    //if (!session) {
    //    return Response.json({ "error": "Unauthorized" }, { status: 401 })
    //}


    try {
        const form = await req.formData()

        const result = await fetch(API_URL + "/tasks/remove", {
            method: "POST",
            body: form,
            headers: {
                "Authorization" : process.env.KITSUNEC2_API_AUTH_TOKEN
            }
        })

        if (result.status === 200) {
            return Response.json({ "success": true })
        } else {
            const error = await result.json()
            return Response.json(error , { status: 500 })
        }
    } catch (e) {
        return Response.json({ "error": e.message }, { status: 500 })
    }

}