import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function POST(req) {
    const session = await getServerSession(authOptions)
    if (!session) {
        return Response.json({ "error": "Unauthorized" }, { status: 401 })
    }


    try {
        const body = await req.json()
        const network = body.network
        const port = body.port
        const params = new URLSearchParams();
        if(network != null){
            params.append("network", network)
        }
        params.append("port", port)

        const result = await fetch(API_URL + "/listeners/add", {
            method: "POST",
            headers: {
                "Authorization" : process.env.KITSUNEC2_API_AUTH_TOKEN
            },
            body: params
        })

        if (result.status === 200) {
            return Response.json({ "success": true })
        } else {
            const error = await result.json()
            return Response.json({ "error": error }, { status: 500 })
        }
    } catch (e) {
        return Response.json({ "error": e.message }, { status: 500 })
    }

}