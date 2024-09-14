import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"
import { Validation } from "@/constants/validation"
import validate from "validate.js"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function GET(req){
    const session = await getServerSession(authOptions)
    if (!session){
        return Response.json({"error":"Unauthorized"}, {status: 401})
    } 

    let taskId = req.nextUrl.searchParams.get("taskId")
    let validationFailed = validate({taskId}, Validation.api_kitsune_file_download)
    if (validationFailed){
        return Response.json({"error" : validationFailed}, { status: 400 })
    }

    try{
        const response = await fetch(API_URL + "/file/download?taskId=" + taskId, {
            cache: "no-store",
            method: 'GET',
            headers: {
                "Authorization" : process.env.KITSUNEC2_API_AUTH_TOKEN
            }
        })
        if (response.status === 200) {
            const blob = await response.blob()
            return new Response(blob)
        } else {
            const error = await response.json()
            return Response.json(error, { status: 500 })
        }
    } catch(e){
        return Response.json({"error":e.message}, {status:500})
    }   
}