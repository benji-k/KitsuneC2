import GenImplantForm from "@/components/genImplantForm"

export default function Generate(){

    return(
        <>
        <h2 className="text-white text-3xl pl-5 pt-5">Generate New Implant Binary</h2>
        <div className="m-5 mt-10">
            <GenImplantForm />
        </div>
        </>
    )
}