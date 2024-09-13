import ListenerTable from "@/components/tables/listenerTable"
import NewListenerForm from "@/components/forms/newListenerForm"

export default function Listeners(){
    return (
    <>
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