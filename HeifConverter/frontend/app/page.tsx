"use client";
import React, { useCallback, useState } from "react";
import { useDropzone } from "react-dropzone";
import { OpenFilePicker, ConvertFile } from "../wailsjs/wailsjs/go/main/App";
import prettyMilliseconds from "pretty-ms";

type FileConversionInfo = {
  OriginalFileName: string;
  OriginalFileSize: number;
  NewFileName: string;
  NewFileSize: number;
  ConversionTime: number;
  Thumbnail: number[];
};

export default function Home() {
  const [state, setState] = useState({
    progress: 0,
    isCompleted: false,
    totalConversionTime: 0,
  });
  const [convertedFiles, setConvertedFiles] = useState<FileConversionInfo[]>(
    []
  );
  const [previewImage, setPreviewImage] = useState<string | null>(null);

  const convertFiles = async (files: File[]) => {
    const startTime = Date.now();
    for (let i = 0; i < files.length; i++) {
      const fileInfo: FileConversionInfo = await ConvertFile(files[i].name);

      setConvertedFiles((prevFiles) => [...prevFiles, fileInfo]);
      setState((prevState) => ({
        ...prevState,
        progress: prevState.progress + 100 / files.length,
        isCompleted: i === files.length - 1,
      }));
    }
    const totalConversionTime = Date.now() - startTime;
    setState((prevState) => ({
      ...prevState,
      totalConversionTime,
    }));
  };

  const onDrop = useCallback((acceptedFiles: File[]) => {
    setState({ progress: 0, isCompleted: false, totalConversionTime: 0 });
    setConvertedFiles([]);
    convertFiles(acceptedFiles);
  }, []);

  const handleFilePicker = async () => {
    const files = await OpenFilePicker();
    setState({ progress: 0, isCompleted: false, totalConversionTime: 0 });
    setConvertedFiles([]);
    const fileObjects = files.map((fileName) => new File([], fileName));
    convertFiles(fileObjects);
  };

  const { getRootProps, getInputProps } = useDropzone({
    onDrop,
    multiple: true,
    onDragEnter: () => {},
    onDragOver: () => {},
    onDragLeave: () => {},
  });

  const formatSize = (size: number) => {
    const i = Math.floor(Math.log(size) / Math.log(1024));
    return (
      Number((size / Math.pow(1024, i)).toFixed(2)) +
      " " +
      ["B", "KB", "MB", "GB", "TB"][i]
    );
  };

  const handlePreview = (filePath: string) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      setPreviewImage(e.target?.result as string);
    };
    reader.readAsDataURL(new File([filePath], filePath));
  };

  const closeModal = () => {
    setPreviewImage(null);
  };

  return (
    <div>
      <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 font-sans">
        <div
          {...getRootProps()}
          className="flex flex-col items-center justify-center w-80 h-48 border-2 border-dashed border-gray-400 bg-white text-gray-700 p-5 mb-5 rounded-lg cursor-pointer"
        >
          <input
            {...getInputProps()}
            type="file"
            accept=".heic,.heif"
            multiple
          />
          <p>
            Drag &apos;n&apos; drop some files here, or click to select files
          </p>
        </div>
        <button
          onClick={handleFilePicker}
          className="px-4 py-2 mb-5 text-white bg-blue-500 rounded hover:bg-blue-600"
        >
          Open File Picker
        </button>
        {!state.isCompleted && state.progress > 0 && (
          <div className="w-full max-w-md mb-5">
            <div className="relative h-8 bg-gray-300 rounded">
              <div
                className="absolute top-0 left-0 h-full bg-blue-500 rounded"
                style={{ width: `${state.progress}%` }}
              >
                <div className="flex items-center justify-center h-full text-white">
                  {state.progress > 0 ? Math.round(state.progress) + "%" : ""}
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
                <th className="px-4 py-2">Time to convert</th>
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
                    {" "}
                    <button
                      className="text-blue-500 underline"
                      onClick={() => handlePreview(file.NewFileName)}
                    >
                      {file.NewFileName}
                    </button>
                  </td>
                  <td className="px-4 py-2">{formatSize(file.NewFileSize)}</td>
                  <td className="px-4 py-2">
                    {prettyMilliseconds(file.ConversionTime)}
                  </td>
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
