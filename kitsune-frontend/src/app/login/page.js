import Image from 'next/image'

export default function Login() {

    return (
        <div className="w-screen h-screen bg-gradient-to-br from-wine-purple from-10% to-zinc-600">
            <div className="flex flex-col items-center justify-start gap-3 mx-auto w-[700px]">
                <Image src="/fox.png" width={141} height={141} alt='Kitsune logo' className='mt-[224px] mb-[136px]'></Image>
                <input type="text" placeholder="Username" className="text-white text-opacity-70 w-[507px] ml-[95px] text-xl font-normal font-['Inter'] bg-transparent outline-none self-start"></input>
                <div className="w-[507px] h-[0px] mb-[20px] border border-white border-opacity-40"></div>
                <input type="password" placeholder="Password" className="text-white text-opacity-70 w-[507px] text-xl ml-[95px] font-normal font-['Inter'] bg-transparent outline-none self-start"></input>
                <div className="w-[507px] h-[0px] border border-white border-opacity-40"></div>
                <button className="w-[543px] h-[67px] mt-[30px] bg-white bg-opacity-90 rounded-[15px] text-neutral-800 text-[32px] font-normal font-['Inter']">Login</button>
            </div>
        </div>
    )
}