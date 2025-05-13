// page.js
"use client";

import Head from "next/head";
import Link from "next/link";
import { useEffect, useState } from "react";
import Image from "next/image";
import { motion, AnimatePresence } from 'framer-motion';
import { Combobox } from "@headlessui/react";
import { CheckIcon, ChevronUpDownIcon } from "@heroicons/react/20/solid";
// import data from './recipe/recipe.json';
import React from "react";
import { ChevronRightIcon, ChevronDownIcon } from "@heroicons/react/20/solid";
import TreeChart from "./TreeChart";


export default function Home() {
  const [data, setData] = useState([]);
  const [mode, setMode] = useState('');
  const [recipeName, setRecipeName] = useState('');
  const [isMultiple, setIsMultiple] = useState(false);
  const [maxRecipes, setMaxRecipes] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [treeData, setTreeData] = useState(null);
  const [nodesVisited, setNodesVisited] = useState(0);
  const [searchTime, setSearchTime] = useState("");
  const [query, setQuery] = useState("");
  const [collapsedNodes, setCollapsedNodes] = useState({});
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/recipe.json")
      .then((res) => res.json())
      .then((json) => setData(json))
      .catch((err) => console.error("Fetch error:", err));
  }, []);

  const elements = (data || []).map(item => item.element);
  const filteredElements =
    query === ""
      ? []
      : elements.filter((element) =>
          element.toLowerCase().includes(query.toLowerCase())
        );
  // console.log(elements); // Ini array berisi semua nama elemen

  const toggleCollapse = (name) => {
    setCollapsedNodes(prev => ({
      ...prev,
      [name]: !prev[name]
    }));
  };

  const renderTree = (node, depth = 0) => {
    if (!node) return null;
  
    const isCollapsed = collapsedNodes[node.element] ?? false;
    const hasChildren = node.children && node.children.length > 0;
  
    const svgPath = `/path/to/svgs/${node.element.toLowerCase()}.svg`; // Misalnya file SVG ada di folder public/svgs/
  
    return (
      <motion.div
        initial={{ opacity: 0, y: 5 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3, delay: depth * 0.05 }}
        className="relative pl-6 border-l-2 border-gray-300 ml-2"
      >
        <div
          className="flex items-center cursor-pointer mb-1 text-black font-medium relative"
          onClick={() => hasChildren && toggleCollapse(node.element)}
        >
          {hasChildren && (
            isCollapsed ? (
              <ChevronRightIcon className="h-4 w-4 text-gray-500 mr-1" />
            ) : (
              <ChevronDownIcon className="h-4 w-4 text-gray-500 mr-1" />
            )
          )}
          <div className="w-2 h-2 bg-blue-500 rounded-full absolute -left-3 top-1" />
          
          {/* Gambar SVG */}
          <Image src={svgPath} alt={node.element} width={20} height={20} />
          
          <span className="ml-2">{node.element}</span>
        </div>
        {/* Keterangan */}
        <p className="text-sm text-gray-600">{node.recipe}</p>
  
        <AnimatePresence>
          {!isCollapsed && hasChildren && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.3 }}
              className="ml-2 overflow-hidden"
            >
              {node.children.map((child, idx) => (
                <div key={idx}>{renderTree(child, depth + 1)}</div>
              ))}
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    );
  };

  const handleSearch = async () => {
    console.log("Sending request with values:");
    console.log("Mode:", mode);
    console.log("Recipe Name:", recipeName);
    console.log("Multiple:", isMultiple);
    console.log("Max Recipes:", maxRecipes);
  
    if (!mode || !recipeName.trim() || (isMultiple && (!maxRecipes || isNaN(maxRecipes)))) {
      setErrorMessage('Please complete all required fields before searching.');
      return;
    }
  
    setIsLoading(true);
    setErrorMessage('');
    setTreeData([]);
    setSearchTime('');
    setNodesVisited(0);
  
    const start = performance.now();
  
    try {
      const requestBody = {
        method: mode,
        target: recipeName,
        multiple: isMultiple,
        ...(isMultiple && { maxRecipes: parseInt(maxRecipes) || 1 })
      };
  
      console.log("Request Body:", requestBody);
  
      const res = await fetch("http://localhost:8080/api/recipe", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestBody),
      });
  
      const rawText = await res.text();
      console.log("Response Status:", res.status);
      console.log("Raw Response Text:", rawText);
  
      if (!res.ok) {
        throw new Error(`Server error (${res.status}): ${rawText}`);
      }
  
      const data = JSON.parse(rawText);
      console.log("Parsed JSON:", data);
  
      setTreeData(data.trees || []);
      setNodesVisited(data.nodesVisited);
      setSearchTime(((performance.now() - start) / 1000).toFixed(2) + "s");
    } catch (error) {
      console.error("Search failed:", error);
      setErrorMessage(error.message || "Search failed.");
    } finally {
      setIsLoading(false); // ⬅️ Selesai loading
    }
  };  
  
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
            <Link href='/aboutpage'>
              <button className='bg-[#D5D5D5] text-black py-2 px-4 rounded h-9'>About Us</button>
            </Link>
          </div>
        </div>

        <div className="pl-2 pr-2 flex flex-row flex-grow overflow-hidden">
          <div className="flex flex-col flex-grow mr-2 rounded-lg w-1/4 overflow-hidden pb-2">
            <div className="bg-[#D5D5D5] p-4 w-[100%] h-[35%] flex-grow overflow-hidden space-y-2">
              <div className="flex flex-col items-center justify-center space-x-2">

                {/* First Question */}
                <p className="text-black text-center font-bold mr-2">Which method do you want to use?</p>
                  {/* Toggle Mode */}
                  <div className="flex justify-center space-x-4 pt-4">
                    <button
                      className={`px-4 py-2 rounded-lg font-semibold ${
                        mode === 'bfs' ? 'bg-[#F9A71F] text-white' : 'bg-gray-200'
                      }`}
                      onClick={() => setMode('bfs')}
                    >
                      BFS
                    </button>
                    <button
                      className={`px-4 py-2 rounded-lg font-semibold ${
                        mode === 'dfs' ? 'bg-[#F9A71F]  text-white' : 'bg-gray-200'
                      }`}
                      onClick={() => setMode('dfs')}
                    >
                     DFS
                    </button>
                  </div>
              </div>

              {/* Second Question */}
              <div className="flex flex-col items-center justify-center space-x-2">
                <p className="text-black text-center font-bold mr-2 pt-4 pb-4">What recipe do you want?</p>
                <Combobox value={recipeName} onChange={setRecipeName}>
                  <div className="relative w-full">
                    <div className="relative w-full cursor-default overflow-hidden rounded bg-white text-left shadow-md focus:outline-none sm:text-sm">
                      <Combobox.Input
                        className="w-full border-none py-2 pl-3 pr-10 text-center leading-5 text-gray-900 focus:ring-0"
                        displayValue={(element) => element}
                        onChange={(event) => setQuery(event.target.value)}
                        placeholder="Enter the recipe name here"
                      />
                      <Combobox.Button className="absolute inset-y-0 right-0 flex items-center pr-2">
                        <ChevronUpDownIcon className="h-5 w-5 text-gray-400" />
                      </Combobox.Button>
                    </div>
                    {filteredElements.length > 0 && (
                      <Combobox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm z-10">
                        {filteredElements.map((element, idx) => (
                          <Combobox.Option
                            key={idx}
                            value={element}
                            className={({ active }) =>
                              `relative cursor-default select-none py-2 pl-10 pr-4 ${
                                active ? 'bg-blue-600 text-white' : 'text-gray-900'
                              }`
                            }
                          >
                            {({ selected, active }) => (
                              <>
                                <span
                                  className={`block truncate ${
                                    selected ? 'font-medium' : 'font-normal'
                                  }`}
                                >
                                  {element}
                                </span>
                                {selected ? (
                                  <span
                                    className={`absolute inset-y-0 left-0 flex items-center pl-3 ${
                                      active ? 'text-white' : 'text-blue-600'
                                    }`}
                                  >
                                    <CheckIcon className="h-5 w-5" aria-hidden="true" />
                                  </span>
                                ) : null}
                              </>
                            )}
                          </Combobox.Option>
                        ))}
                      </Combobox.Options>
                    )}
                  </div>
                </Combobox>

             {/* Third Question */}
              <div className="flex flex-col items-center justify-center space-x-2">
                <p className="text-black text-center font-bold mr-2 pt-4">Find one or many recipes?</p>

                <div className="space-y-6 text-center pt-4 flex flex-col items-center">
                    {/* Toggle Title */}
                    <h1 className="text-l font-bold">
                      {isMultiple ? 'Multiple Recipes' : 'One Recipe'}
                    </h1>

                    {/* Toggle Switch */}
                    <div
                      onClick={() => setIsMultiple(!isMultiple)}
                      className={`w-20 h-10 rounded-full flex items-center px-1 cursor-pointer transition-colors duration-300 ${
                        isMultiple ? 'bg-purple-800' : 'bg-purple-700'
                      }`}
                    >
                      <motion.div
                        layout
                        transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                        className={`w-8 h-8 rounded-full bg-white shadow-md ${
                          isMultiple ? 'ml-0' : 'ml-auto'
                        }`}
                      />
                    </div>

                    {/* Text content animation */}
                    <AnimatePresence mode="wait">
                      {isMultiple ? (
                        <motion.p
                          key="multiple"
                          initial={{ opacity: 0, y: 10 }}
                          animate={{ opacity: 1, y: 0 }}
                          exit={{ opacity: 0, y: -10 }}
                          transition={{ duration: 0.3 }}
                          className="text-gray-700"
                        >
                          You will get multiple different recipes for the target element.
                        </motion.p>
                      ) : (
                        <motion.p
                          key="one"
                          initial={{ opacity: 0, y: 10 }}
                          animate={{ opacity: 1, y: 0 }}
                          exit={{ opacity: 0, y: -10 }}
                          transition={{ duration: 0.3 }}
                          className="text-gray-700"
                        >
                          You will get only one recipe — any valid one is okay.
                        </motion.p>
                      )}
                    </AnimatePresence>

                    {/* Max recipes input (only if multiple) */}
                    {isMultiple && (
                      <div className="text-center">
                        <label className="block mb-2 font-medium">
                          How many recipes do you want to find at most?
                        </label>
                        <input
                          type="number"
                          min={1}
                          max={720}
                          value={maxRecipes === '' ? '' : String(maxRecipes)}
                          onChange={(e) => {
                            const value = parseInt(e.target.value);
                            if (!isNaN(value)) {
                              setMaxRecipes(Math.min(Math.max(value, 1), 720)); // Clamp antara 1 dan 720
                            } else {
                              setMaxRecipes('');
                            }
                          }}
                          className="border rounded px-3 py-1 w-48 text-center"
                        />
                      </div>
                    )}
                  </div>
              </div>

              {/* Button Search */}
              <div className="flex flex-col items-center justify-center mt-6">
                {errorMessage && (
                  <p className="text-red-600 mb-2 font-medium text-center">
                    {errorMessage}
                  </p>
                )}
                <button
                  className="bg-[#451952] hover:bg-[#5e2470] text-white py-2 px-6 rounded transition-colors duration-200"
                  onClick={handleSearch}
                  disabled={false}
                >
                  Search
                </button>
              </div>
            </div>
          </div>
        </div>

        <div className="flex flex-col flex-grow rounded-lg w-3/4 overflow-hidden pb-2">
          <div className="bg-[#D5D5D5] p-4 w-full h-full flex-grow overflow-hidden space-y-2">
            {/* Tree Display Area */}
            <div className="border rounded-lg bg-white flex flex-col justify-between h-[639px] p-4">
            {isLoading ? (
              <div className="flex justify-center items-center flex-grow">
                <div className="text-center">
                  <div className="loader mx-auto mb-4"></div>
                  <p className="text-gray-600">Loading search results...</p>
                </div>
              </div>
            ) : treeData && treeData.length > 0 ? (
              <div className="flex flex-col bg-gray-50 p-4 space-y-4 flex-grow overflow-auto">
                <p className="text-black text-center font-bold">Search Results</p>
                <div className="space-y-4">
                  {treeData.map((tree, index) => (
                    <TreeChart key={index} data={tree} />
                  ))}
                </div>
              </div>
            ) : (
              <div className="flex justify-center items-center flex-grow">
                <p className="text-gray-500">No results found. Please search first!</p>
              </div>
            )}

              {/* Bottom Info */}
              <div className="flex justify-between pt-4 text-sm text-black font-medium">
                <div className="ml-2">Search Time: <span className="font-normal">{searchTime || "-"}</span></div>
                <div className="mr-2">Number of Nodes Visited: <span className="font-normal">{nodesVisited || "-"}</span></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
   </>
  );
}