import { Download, Trash2 } from "lucide-react"

function formatSize(bytes) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString("en-US", {
    month: "short", day: "numeric", year: "numeric"
  })
}

export default function FileList({ files, onDelete }) {
  if (files.length === 0) {
    return <p className="empty">No files yet. Upload something!</p>
  }

  return (
    <ul className="file-list">
      {files.map(file => (
        <li key={file.name} className="file-item">
          <div className="file-info">
            <span className="file-name">{file.name}</span>
            <span className="file-meta">
              {formatSize(file.size)} · {formatDate(file.modified_at)}
            </span>
          </div>
          <div className="file-actions">
            <a
              href={`/api/download/${file.name}`}
              download={file.name}
              className="icon-btn download"
              title="Download"
            >
              <Download size={16} />
            </a>
            <button
              onClick={() => onDelete(file.name)}
              className="icon-btn delete"
              title="Delete"
            >
              <Trash2 size={16} />
            </button>
          </div>
        </li>
      ))}
    </ul>
  )
}