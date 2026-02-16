import { useState, useEffect } from 'react'
import { useWebSocket } from './hooks/useWebSocket'
import { Sidebar } from './components/Sidebar'
import { ChatArea } from './components/ChatArea'
import { Header } from './components/Header'
import type { Session, Message } from './types'

function App() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [currentSession, setCurrentSession] = useState<string | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const [isSidebarOpen, setIsSidebarOpen] = useState(true)
  
  const { sendMessage, lastMessage, connectionStatus } = useWebSocket('ws://localhost:8080/ws')

  useEffect(() => {
    setIsConnected(connectionStatus === 'connected')
  }, [connectionStatus])

  useEffect(() => {
    if (lastMessage) {
      handleWebSocketMessage(lastMessage)
    }
  }, [lastMessage])

  const handleWebSocketMessage = (data: any) => {
    if (data.type === 'event') {
      switch (data.event) {
        case 'connected':
          console.log('Connected to HiveClaw gateway')
          // Request session list
          sendMessage({ type: 'req', id: '1', method: 'session.list', params: {} })
          break
        case 'message':
          if (data.payload.sessionId === currentSession) {
            setMessages(prev => [...prev, data.payload.message])
          }
          break
      }
    } else if (data.type === 'res') {
      if (data.ok && data.payload) {
        // Handle session list response
        if (Array.isArray(data.payload)) {
          setSessions(data.payload)
        }
        // Handle chat response
        if (data.payload.response) {
          setMessages(prev => [...prev, {
            id: Date.now().toString(),
            role: 'assistant',
            content: data.payload.response,
            timestamp: new Date().toISOString()
          }])
        }
      }
    }
  }

  const handleNewSession = () => {
    const id = `session_${Date.now()}`
    sendMessage({
      type: 'req',
      id: '2',
      method: 'session.create',
      params: { name: `New Chat ${sessions.length + 1}` }
    })
    setCurrentSession(id)
    setMessages([])
  }

  const handleSelectSession = (sessionId: string) => {
    setCurrentSession(sessionId)
    const session = sessions.find(s => s.id === sessionId)
    if (session) {
      setMessages(session.messages || [])
    }
  }

  const handleSendMessage = (content: string) => {
    if (!content.trim()) return

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content,
      timestamp: new Date().toISOString()
    }

    setMessages(prev => [...prev, userMessage])

    sendMessage({
      type: 'req',
      id: Date.now().toString(),
      method: 'chat.send',
      params: {
        sessionId: currentSession || 'main',
        message: content
      }
    })
  }

  return (
    <div className="flex h-screen bg-dark-950">
      {/* Sidebar */}
      <Sidebar
        sessions={sessions}
        currentSession={currentSession}
        isOpen={isSidebarOpen}
        onToggle={() => setIsSidebarOpen(!isSidebarOpen)}
        onNewSession={handleNewSession}
        onSelectSession={handleSelectSession}
      />

      {/* Main Content */}
      <div className="flex-1 flex flex-col">
        <Header 
          isConnected={isConnected}
          onToggleSidebar={() => setIsSidebarOpen(!isSidebarOpen)}
        />
        
        <ChatArea
          messages={messages}
          onSendMessage={handleSendMessage}
          currentSession={currentSession}
        />
      </div>
    </div>
  )
}

export default App
