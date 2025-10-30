function CountersBanner() {
  return (
    <div className="flex justify-center group mb-15 xl:mb-20">
      <div className="relative flex justify-center items-center max-w-[978px] w-full overflow-hidden min-h-21 border border-line mx-5 group-hover:border-gray-300">
        <a className="group-hover:opacity-65 duration-300" href="#">
          <img className="absolute w-50 h-[105px] left-0 top-0 translate-x-[-40%] -translate-y-[-5%]
        sm:translate-x-[-35%] md:translate-x-[-25%] lg:translate-x-[-15%] xl:translate-x-[-5%]" src="counter.svg" alt="counter" />
          <p className="relative tracking-[1px] text-[21px] text-center uppercase bg-transparent z-1 w-[286px] avg:w-auto">
            Пора передавать показания счетчиков
          </p>
          <img className="absolute w-50 h-[105px] right-0 top-0 -translate-y-[30%] -translate-x-[-40%]
        sm:-translate-x-[-35%] md:-translate-x-[-25%] lg:-translate-x-[-15%] xl:-translate-x-[-5%]" src="counter.svg" alt="counter" />
        </a>
      </div>
    </div >
  );
}

export default CountersBanner;