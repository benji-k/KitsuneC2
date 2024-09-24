import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"
import { Validation } from "@/constants/validation"
import validate from "validate.js"

const API_URL = process.env.KITSUNEC2_WEBURL

export async function POST(req) {
    const session = await getServerSession(authOptions)
    if (!session) {
        return Response.json({ "error": "Unauthorized" }, { status: 401 })
    }

    const validateParams = function (form) {
        const data = {};
        form.forEach((value, key) => {
            data[key] = value;
        });

        return validate(data, Validation.api_kitsune_tasks_add)
    }

    const form = await req.formData()

    const validationFailed = validateParams(form)
    if (validationFailed) {
        return Response.json({ "error": validationFailed }, { status: 400 })
    }

    const implants = form.get("implants")
    form.set("implants", "[" + implants + "]") //the backend expects an array [] 

    try {
        const result = await fetch(API_URL + "/tasks/add", {
            method: "POST",
            body: form,
            headers: {
                "Authorization": process.env.KITSUNEC2_API_AUTH_TOKEN
            }
        })

        if (result.status === 200) {
            return Response.json({ "success": true })
        } else {
            const error = await result.json()
            return Response.json(error, { status: 500 })
        }
    } catch (e) {
        return Response.json({ "error": e.message }, { status: 500 })
    }
}