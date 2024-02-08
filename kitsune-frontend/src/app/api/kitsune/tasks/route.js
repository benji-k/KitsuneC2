import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function GET(req){
    const session = await getServerSession(authOptions)
    if (!session){
        return Response.json({"error":"Unauthorized"}, {status: 401})
    } 

    let completed = req.nextUrl.searchParams.get("completed")
    try{
        const response = await fetch(API_URL + (completed == "true" ? "/tasks?completed=true" : "/tasks?completed=false"), {cache: "no-store"})
        const responseJson = await response.json()
        if (response.status == 200){
            return Response.json(responseJson)
        } else {
            return Response.json(responseJson, {status:500})
        }
        
    } catch{
        return Response.json({"error":"Kitsune server unreachable"}, {status:500})
    }   
}