import type { EventDetails, EventPayload, EventSummary, Profile, Session, TicketType } from '../types'

const API_BASE_URL = 'http://localhost:8080/v1/api'
const SESSION_KEY = 'ms-congress-session'
const COMPANY_TYPE = 2

export class ApiError extends Error {}

async function request<T>(path: string, options: RequestInit = {}, session?: Session): Promise<T> {
  const headers = new Headers(options.headers)
  if (options.body) headers.set('Content-Type', 'application/json')
  if (session?.token) headers.set('Authorization', `${session.tokenType ?? 'Bearer'} ${session.token}`)

  const response = await fetch(`${API_BASE_URL}${path}`, { ...options, headers })
  if (!response.ok) {
    const raw = await response.text()
    let data: unknown
    try { data = raw ? JSON.parse(raw) : null } catch { throw new ApiError(raw || `A API respondeu com status ${response.status}.`) }
    const body = data as { erro?: string; message?: string; detalhes?: string } | null
    throw new ApiError(body?.erro ?? body?.message ?? body?.detalhes ?? `A API respondeu com status ${response.status}.`)
  }
  const raw = await response.text()
  return (raw ? JSON.parse(raw) : null) as T
}

export function getStoredSession(): Session | null {
  const saved = localStorage.getItem(SESSION_KEY)
  if (!saved) return null
  try { return JSON.parse(saved) as Session } catch { return null }
}

export function saveSession(session: Session) { localStorage.setItem(SESSION_KEY, JSON.stringify(session)) }
export function clearSession() { localStorage.removeItem(SESSION_KEY) }

export async function login(email: string, senha: string): Promise<Session> {
  const session = await request<Session>('/login', { method: 'POST', body: JSON.stringify({ email, senha }) })
  if (Number(session.tipoPessoaID) !== COMPANY_TYPE) throw new ApiError('Este painel está disponível apenas para contas de pessoa jurídica.')
  return session
}

export function getProfile(session: Session) { return request<Profile>(`/usuarios/${session.id}`, {}, session) }
export function getEvents(session: Session, organizerId: number) { return request<EventSummary[]>(`/eventos/organizador/${organizerId}`, {}, session) }
export function getTicketTypes(session: Session) { return request<TicketType[]>('/eventos/tipos-ingresso', {}, session) }
export function getEvent(session: Session, id: number) { return request<EventDetails>(`/eventos/${id}`, {}, session) }
export function createEvent(session: Session, payload: EventPayload) { return request<EventDetails>('/eventos', { method: 'POST', body: JSON.stringify(payload) }, session) }
export function updateEvent(session: Session, id: number, payload: EventPayload) { return request<EventDetails>(`/eventos/${id}`, { method: 'PUT', body: JSON.stringify(payload) }, session) }
export function publishEvent(session: Session, id: number) { return request<EventDetails>(`/eventos/${id}/publicar`, { method: 'PATCH' }, session) }
export function closeEvent(session: Session, id: number) { return request<EventDetails>(`/eventos/${id}/encerrar`, { method: 'PATCH' }, session) }
export function deleteEvent(session: Session, id: number) { return request<void>(`/eventos/${id}`, { method: 'DELETE' }, session) }
