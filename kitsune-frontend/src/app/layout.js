import './globals.css'
import SessionProvider from '@/components/sessionProvider'
import { getServerSession } from 'next-auth'

export default async function RootLayout({ children }) {
const session = await getServerSession()

  return (
    <html lang="en">
      <SessionProvider session={session}>
        <body>{children}</body>
      </SessionProvider>
    </html>
  )
}
