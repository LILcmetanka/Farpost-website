import { useState, type PropsWithChildren } from "react";
import Header from "../header/Header";
import Footer from "../footer/Footer";
import Message from "../message/Message";
import CountersBanner from "../counters-info/CountersBanner";

function Layout({ children }: PropsWithChildren) {
	const [isMessageVisible, setIsMessageVisible] = useState(true);
	
	return (
		<div className="text-[14px] font-sans">
			<Message isClick={isMessageVisible} setIsClick={setIsMessageVisible}/>
			<Header isMessageVisible={isMessageVisible}/>
			<CountersBanner />
			<main>
				{children}
			</main>
			<hr className="my-10 h-px bg-line border-none"/>
			<Footer />
		</div>
	);
}

export default Layout;