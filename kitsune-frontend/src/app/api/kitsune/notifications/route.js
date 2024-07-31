import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function GET(req){
    const session = await getServerSession(authOptions)
    if (!session){
        return Response.json({"error":"Unauthorized"}, {status: 401})
    }


    try{
        const response = await fetch(API_URL + "/notifications", {
            cache: "no-store",
            headers: {
                "Authorization" : process.env.KITSUNEC2_API_AUTH_TOKEN
            }
        })
        const responseJson = await response.json()
        
        //we add timestamps to the notifications because of issues with SWR not re-rendering components if the data stays the same
        const responseWithTimestamps = responseJson.map(notification => ({
            ...notification,
            Timestamp: Date.now()
          }));

        if (response.status == 200){
            return Response.json(responseWithTimestamps)
        } else {
            return Response.json(responseWithTimestamps, {status:500})
        }
        
    } catch(error){
        return Response.json({"error":"Kitsune server unreachable"}, {status:500})
    }   

}