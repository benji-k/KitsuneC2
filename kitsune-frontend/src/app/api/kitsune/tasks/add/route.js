import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"
import { Tasks } from "@/constants/tasks"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function POST(req){
    const session = await getServerSession(authOptions)
    if (!session){
        return Response.json({"error":"Unauthorized"}, {status: 401})
    } 

    const body = await req.json()
    
    const implants = body.implants
    const taskType = body.taskType
    const taskArguments = body.arguments

    if(implants.length === 0){
        return Response.json({"error" : "At least one implant should be selected"}, {status:500})
    }

    if (!Tasks.find((t) => (t.taskType === taskType))){
        return Response.json({"error" : "invalid task type"}, {status:500})
    }

    try{
        const params = new URLSearchParams();

        params.append("implants", "[" + implants + "]")
        params.append("taskType", taskType)
        Object.keys(taskArguments).forEach(key => {
            params.append(key, taskArguments[key]);
        });

        const result = await fetch(API_URL + "/tasks/add", {
            method: "POST",
            body: params
        })

        if (result.status === 200){
            return Response.json({"success" : true})
        } else {
            const error = await result.json()
            return Response.json({"error":error}, {status:500})
        }
    } catch(e){
        return Response.json({"error":e.message}, {status:500})
    }
}