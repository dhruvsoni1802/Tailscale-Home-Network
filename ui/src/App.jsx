import { useState, useEffect, useCallback } from "react"
import FileList from "./components/FileList"
import UploadButton from "./components/UploadButton"
import { HardDrive } from "lucide-react"
import "./App.css"

export default function App() {
  const [files, setFiles] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchFiles = useCallback(async () => {
    try {
      const resp = await fetch("/api/files")
      if (!resp.ok) throw new Error("failed to fetch files")
      const data = await resp.json()
      setFiles(data.files || [])
      setError(null)
    } catch (err) {
      setError("Could not reach storage server.")
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetchFiles() }, [fetchFiles])

  async function handleDelete(filename) {
    if (!confirm(`Delete ${filename}?`)) return

    const resp = await fetch(`/api/files/${filename}`, { method: "DELETE" })
    if (resp.ok) fetchFiles()
    else alert("Delete failed")
  }

  return (
    <div className="app">
      <header className="header">
        <div className="header-inner">
          <div className="logo">
            <HardDrive size={22} />
            <span>TailStore</span>
          </div>
          <UploadButton onUploaded={fetchFiles} />
        </div>
      </header>

      <main className="main">
        <div className="card">
          <div className="card-header">
            <h2>Your Files</h2>
            <span className="count">{files.length} file{files.length !== 1 ? "s" : ""}</span>
          </div>

          {loading && <p className="empty">Loading...</p>}
          {error && <p className="error">{error}</p>}
          {!loading && !error && (
            <FileList files={files} onDelete={handleDelete} />
          )}
        </div>
      </main>
    </div>
  )
}