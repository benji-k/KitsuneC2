import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function GET(req){
    const session = await getServerSession(authOptions)
    if (!session){
        return Response.json({"error":"Unauthorized"}, {status: 401})
    } 

    try{
        const response = await fetch(API_URL + "/implants", {
            cache: "no-store",
            headers: {
                "Authorization" : process.env.KITSUNEC2_API_AUTH_TOKEN
            },
        })
        const responseJson = await response.json()
        if (response.status == 200){
            return Response.json(responseJson)
        } else {
            return Response.json(responseJson, {status:500})
        }
        
    } catch(error){
        console.log(error)
        return Response.json({"error":error.message}, {status:500})
    }   
}