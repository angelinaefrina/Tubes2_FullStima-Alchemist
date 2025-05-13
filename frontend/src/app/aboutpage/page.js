"use client";

import Head from "next/head";
import Link from "next/link";

export default function Home() {
  return (
    <>
      <Head>
        <title>FullStima Alchemist</title>
        </Head>
      <div className="flex flex-col h-screen overflow-hidden">
        <div className="ml-4 mr-2 mt-1 justify-center h-[7%]">
          <div className="flex flex-row items-center justify-between">
            <div className="flex flex-row items-center space-x-4">
            <img
                  src="/assets/logo.svg"
                  className=" object-cover mb- w-10 h-12 pb-1"
                  alt="Logo"
                />
              <p className="text-[#FFFFFF] text-2xl"><b>FullStima Alchemist</b></p>
              <p className="text-[#FFFFFF] text-m">- Little Alchemy 2 Element Finder</p>
            </div>
            <Link href='/'>
              <button className='bg-[#D5D5D5] text-black py-2 px-4 rounded h-9'>Home</button>
            </Link>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex flex-row flex-grow overflow-y-auto">
          <div className="flex flex-col flex-grow items-center overflow-y-auto">
            <div className="bg-[#d5d5d5] p-4 w-full max-w-3xl h-auto flex-grow space-y-4 rounded-lg shadow-lg">
              {/* Introduction */}
              <p className="text-[#000000] text-center text-xl"><b>Introduction</b></p>
              <div className="rounded-lg bg-white p-4 flex flex-row items-center justify-center">
                <p className="text-[#000000] text-m">
                FullStima Alchemist is a web-based application designed to help players of Little Alchemy 2
                find the most efficient combinations to create elements. Built with a modern tech stack, the
                frontend is developed using Next.js for a responsive and interactive user experience, while
                the backend is powered by Go (Golang) for high performance and concurrency. The core feature
                of FullStima Alchemist is its intelligent search system, which utilizes Breadth-First Search
                (BFS) and Depth-First Search (DFS) algorithms to explore all possible crafting paths from
                basic elements to a target element. With real-time visualizations and intuitive design,
                FullStima Alchemist simplifies the element discovery process and enhances the fun of
                experimentation in Little Alchemy 2.
                </p>
              </div>

              {/* How to Use */}
              <p className="text-[#000000] text-center text-xl mb-2"><b>How To Use</b></p>
              <div className="rounded-lg bg-white p-4 flex flex-row items-center">
                <p className="text-[#000000] text-m">
                  1. Choose a search method: Click either <strong>BFS</strong> or <strong>DFS</strong> to select the algorithm.<br />
                  2. Enter the target element: Type the name of the element you want to search for.<br />
                  3. Select recipe mode: Toggle between <strong>One Recipe</strong> or <strong>Multiple Recipes</strong>.<br />
                  4. (Optional) If you chose multiple recipes, specify the maximum number of results you want.<br />
                  5. Click <strong>Search</strong> to start finding recipes.
                </p>

              </div>

              {/* Creators */}
              <p className="text-[#000000] text-center text-xl"><b>Creators</b></p>
              <div className="rounded-lg bg-white p-4 flex flex-col items-center">
                <img
                  src="/assets/creators.jpg"
                  className=" object-cover mb-4 w-40 h-40 rounded-lg"
                  alt="Creators"
                />
                <p className="text-center text-[#000000] text-lg">
                  <b>Kelompok 47 - FullStima Alchemist</b><br />
                  1. Shannon Aurellius Anastasya Lie (13523019)<br />
                  2. Angelina Efrina Prahastaputri (13523060)<br />
                  3. Sebastian Hung Yansen (13523070)
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
