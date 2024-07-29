import Navbar from "@/components/navbar"
import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"
import { redirect } from 'next/navigation';
import NotificationBar from "@/components/notificationBar";

export default async function KitsuneLayout({ children }) {
    const session = await getServerSession(authOptions)
    if (!session){
        return redirect("/api/auth/signin", )
    } else {
        return (
            <>
                <Navbar />
                <div className="bg-kc2-dashboard-bg min-h-screen">
                    <NotificationBar popupTime={3000} refreshRate={5000} />
                    {children}
                </div>
            </>
        )
    }
  }