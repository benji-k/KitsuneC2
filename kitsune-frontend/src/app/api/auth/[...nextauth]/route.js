import NextAuth from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"

export const authOptions = {
  pages: {
    signIn: '/login',
  },
  providers: [
      CredentialsProvider({
        name: 'Credentials',
        credentials: {
          username: { label: "Username", type: "text", placeholder: "jsmith" },
          password: { label: "Password", type: "password" }
        },
        async authorize(credentials, req) {
          if (credentials.username === process.env.USERNAME && credentials.password === process.env.PASSWORD){
              return {"authenticated" : true, "username": process.env.USERNAME}
          } else {
              return null
          }
        }
      })
    ]
} 

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST }