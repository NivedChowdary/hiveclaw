import type { Session } from '../types'

interface SidebarProps {
  sessions: Session[]
  currentSession: string | null
  isOpen: boolean
  onToggle: () => void
  onNewSession: () => void
  onSelectSession: (sessionId: string) => void
}

export function Sidebar({
  sessions,
  currentSession,
  isOpen,
  onNewSession,
  onSelectSession,
}: SidebarProps) {
  if (!isOpen) return null

  return (
    <div className="w-64 bg-slate-900 border-r border-slate-700 flex flex-col">
      {/* New Chat Button */}
      <div className="p-3">
        <button
          onClick={onNewSession}
          className="w-full flex items-center gap-2 px-4 py-2.5 bg-amber-600 hover:bg-amber-700 text-white rounded-lg transition-colors font-medium"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          New Chat
        </button>
      </div>

      {/* Sessions List */}
      <div className="flex-1 overflow-y-auto px-2">
        <div className="text-xs font-semibold text-slate-400 px-3 py-2 uppercase tracking-wider">
          Recent Chats
        </div>
        
        {sessions.length === 0 ? (
          <div className="px-3 py-8 text-center text-slate-500">
            <p className="text-sm">No conversations yet</p>
            <p className="text-xs mt-1">Start a new chat to begin</p>
          </div>
        ) : (
          <div className="space-y-1">
            {sessions.map((session) => (
              <button
                key={session.id}
                onClick={() => onSelectSession(session.id)}
                className={`w-full text-left px-3 py-2 rounded-lg transition-colors ${
                  currentSession === session.id
                    ? 'bg-slate-700 text-white'
                    : 'text-slate-300 hover:bg-slate-800 hover:text-white'
                }`}
              >
                <div className="flex items-center gap-2">
                  <svg className="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                  </svg>
                  <span className="truncate text-sm">{session.name}</span>
                </div>
                {session.updatedAt && (
                  <div className="text-xs text-slate-500 mt-0.5 pl-6">
                    {new Date(session.updatedAt).toLocaleDateString()}
                  </div>
                )}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="p-3 border-t border-slate-700">
        <div className="flex items-center gap-3 px-2">
          <div className="w-8 h-8 rounded-full bg-amber-600 flex items-center justify-center text-sm font-medium">
            🐝
          </div>
          <div className="flex-1 min-w-0">
            <div className="text-sm font-medium truncate">HiveClaw</div>
            <div className="text-xs text-slate-400">v0.1.0</div>
          </div>
        </div>
      </div>
    </div>
  )
}
