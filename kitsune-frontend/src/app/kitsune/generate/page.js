import GenImplantForm from "@/components/genImplantForm"
import NotificationBar from "@/components/notificationBar"

export default function Generate(){

    return(
        <>
        <NotificationBar popupTime={4000} />
        <h2 className="text-white text-3xl pl-5 pt-5">Generate New Implant Binary</h2>
        <div className="m-5 mt-10">
            <GenImplantForm />
        </div>
        </>
    )
}