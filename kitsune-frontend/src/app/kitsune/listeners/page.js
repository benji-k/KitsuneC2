import ListenerTable from "@/components/listenerTable"
import NewListenerForm from "@/components/newListenerForm"
import NotificationBar from "@/components/notificationBar"

export default function Listeners(){
    return (
    <>
        <NotificationBar popupTime={4000} />
        <h2 className="text-white text-3xl pl-5 pt-5">Listeners</h2>
        <div className="m-5 mt-3">
            <ListenerTable 
                refreshRate={3000}
            />
        </div>
        <h2 className="text-white text-3xl pl-5 pt-5">Add Listener</h2>
        <div className="m-5 mt-3">
            <NewListenerForm />
        </div>
    </>
    )
}