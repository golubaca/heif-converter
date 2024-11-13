"use client";
import React, { useState, useEffect } from "react";
import { OpenFilePicker, ConvertMultipleFiles } from "../wailsjs/wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/wailsjs/runtime/runtime";
import prettyMilliseconds from "pretty-ms";

type FileConversionInfo = {
  OriginalFileName: string;
  OriginalFileSize: number;
  NewFileName: string;
  NewFileSize: number;
  ConversionTime: number;
  Thumbnail: number[];
};

type ConversionResult = {
  Error: string | null;
  TotalTime: number;
};

export default function Home() {
  const [totalFiles, setTotalFiles] = useState(0);
  const [progress, setProgress] = useState(0);
  const [state, setState] = useState({
    isCompleted: false,
    totalConversionTime: 0,
  });
  const [convertedFiles, setConvertedFiles] = useState<FileConversionInfo[]>(
    []
  );

  useEffect(() => {
    setProgress(Math.round(((convertedFiles.length) * 100) / totalFiles))
    // console.log(`Progress: ${progress}/${totalFiles}`);
  }, [convertedFiles]);

  useEffect(() => {
    EventsOn("conversion_progress", (fileInfo: FileConversionInfo) => {
      setConvertedFiles(prev => [...prev, fileInfo]);
    });

    // Listen for errors
    EventsOn("conversion_error", (error: { path: string, error: string }) => {
      console.error(`Error converting ${error.path}:`, error.error);
    });

    // Listen for completion
    EventsOn("conversion_complete", (results: ConversionResult) => {
      setState(prev => ({
        ...prev,
        isCompleted: true,
        totalConversionTime: results.TotalTime,
      }));
    });
  }, []);

  const convertFiles = async (files: File[]) => {
    setTotalFiles(files.length);
    setState({ isCompleted: false, totalConversionTime: 0});
    setConvertedFiles([]);
    setProgress(0);
    const filePaths = files.map(file => file.name);
    await ConvertMultipleFiles(filePaths);
  };

  const handleFilePicker = async () => {
    const files = await OpenFilePicker();
    setState({ isCompleted: false, totalConversionTime: 0 });
    setConvertedFiles([]);
    const fileObjects = files.map((fileName) => new File([], fileName));
    convertFiles(fileObjects);
  };

  const formatSize = (size: number) => {
    const i = Math.floor(Math.log(size) / Math.log(1024));
    return (
      Number((size / Math.pow(1024, i)).toFixed(2)) +
      " " +
      ["B", "KB", "MB", "GB", "TB"][i]
    );
  };

  return (
    <div>
      <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 font-sans">
        
        <button
          onClick={handleFilePicker}
          className="px-4 py-2 mb-5 text-white bg-blue-500 rounded hover:bg-blue-600"
        >
          Open File Picker
        </button>
        {!state.isCompleted && totalFiles > 0 && (
          <div className="w-full max-w-md mb-5">
            <div className="relative h-8 bg-gray-300 rounded">
              <div
                className="absolute top-0 left-0 h-full bg-blue-500 rounded"
                style={{ width: `${progress}%` }}
              >
                <div className="flex items-center justify-center h-full text-white">
                  {progress > 0 ? progress + "%" : ""}
                </div>
              </div>
            </div>
          </div>
        )}
        {state.isCompleted && (
          <div className="text-lg text-green-500 mb-5">
            Conversion Completed in{" "}
            {prettyMilliseconds(state.totalConversionTime)}!
          </div>
        )}
        {convertedFiles.length > 0 && (
          <table className="table-auto w-full max-w-4xl bg-white rounded-lg shadow-md">
            <thead>
              <tr className="bg-gray-200 text-black">
                <th className="px-4 py-2">Original Name</th>
                <th className="px-4 py-2">Original Size</th>
                <th className="px-4 py-2">New Name</th>
                <th className="px-4 py-2">New Size</th>
                <th className="px-4 py-2">Status</th>
              </tr>
            </thead>
            <tbody>
              {convertedFiles.map((file, index) => (
                <tr
                  key={index}
                  className={`border-t text-black ${
                    index % 2 === 0 ? "bg-gray-100" : "bg-gray-200"
                  }`}
                >
                  <td className="px-4 py-2">{file.OriginalFileName}</td>
                  <td className="px-4 py-2">
                    {formatSize(file.OriginalFileSize)}
                  </td>
                  <td className="px-4 py-2">
                    {file.NewFileName}
                  </td>
                  <td className="px-4 py-2">{formatSize(file.NewFileSize)}</td>
                  <td className="px-4 py-2 text-green-500">✔</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
