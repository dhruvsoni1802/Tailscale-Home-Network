import { useRef, useState } from "react"
import { Upload } from "lucide-react"

export default function UploadButton({ onUploaded }) {
  const inputRef = useRef(null)
  const [uploading, setUploading] = useState(false)

  async function handleFileChange(e) {
    const file = e.target.files[0]
    if (!file) return

    setUploading(true)

    const form = new FormData()
    form.append("file", file)

    try {
      const resp = await fetch("/api/upload", { method: "POST", body: form })
      if (!resp.ok) throw new Error("upload failed")
      onUploaded()
    } catch (err) {
      alert("Upload failed: " + err.message)
    } finally {
      setUploading(false)
      // reset input so same file can be uploaded again
      e.target.value = ""
    }
  }

  return (
    <div>
      <input
        ref={inputRef}
        type="file"
        onChange={handleFileChange}
        style={{ display: "none" }}
      />
      <button
        onClick={() => inputRef.current.click()}
        disabled={uploading}
        className="upload-btn"
      >
        <Upload size={16} />
        {uploading ? "Uploading..." : "Upload File"}
      </button>
    </div>
  )
}