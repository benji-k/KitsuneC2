import { useDashboardState } from "@/state/application"

export default function ConfirmationWindow() {

    const setConfirmationWindowOpen = useDashboardState((state) => state.setConfirmationWindowOpen)
    const confirmationWindowText = useDashboardState((state) => state.confirmationWindowText)
    const confirmationWindowCallbackFn = useDashboardState((state) => state.confirmationWindowCallbackFn)

    return (
        <div className="fixed flex justify-center items-center top-0 left-0 h-full w-full bg-[#9B9B9B]/[.54] z-20">
            <div className="bg-kc2-dark-gray text-white flex flex-col items-center w-full px-3 mx-5 pt-4 rounded-2xl max-w-6xl max-h-[600px] scrollbar-hide overflow-scroll">
               {confirmationWindowText}
                <div className="flex w-full justify-center md:justify-end items-center gap-10 mb-4 mt-4">
                    <button className="bg-[#F96B6B] rounded-md px-8 py-1"
                        onClick={() => { setConfirmationWindowOpen(false) }}>Cancel
                    </button>
                    <button className="bg-[#0EC420] rounded-md px-8 py-1"
                    onClick={() => {
                        confirmationWindowCallbackFn()
                        setConfirmationWindowOpen(false)
                    }}>
                        Confirm
                    </button>
                </div>
            </div>
        </div>
    )
}