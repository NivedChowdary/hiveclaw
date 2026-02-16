export interface Message {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: string
}

export interface Session {
  id: string
  name: string
  agentId: string
  messages: Message[]
  createdAt: string
  updatedAt: string
}

export interface WSMessage {
  type: 'req' | 'res' | 'event'
  id?: string
  method?: string
  event?: string
  params?: any
  payload?: any
  ok?: boolean
  error?: { code: string; message: string }
}
