import type { AvailableEventGroup, Cart, EventDetails, EventPayload, EventSummary, MyEvent, Order, Profile, Session, TicketType } from '../types'

const API_BASE_URL = 'http://localhost:8080/v1/api'
const SESSION_KEY = 'ms-congress-session'

export class ApiError extends Error {}

async function request<T>(path: string, options: RequestInit = {}, session?: Session): Promise<T> {
  const headers = new Headers(options.headers)
  if (options.body) headers.set('Content-Type', 'application/json')
  if (session?.token) headers.set('Authorization', `${session.tokenType ?? 'Bearer'} ${session.token}`)

  const response = await fetch(`${API_BASE_URL}${path}`, { ...options, headers, credentials: 'include' })
  if (!response.ok) {
    if (response.status === 401) clearSession()
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
  return request<Session>('/login', { method: 'POST', body: JSON.stringify({ email, senha }) })
}

export function loginCongressista(email: string, senha: string) {
  return request<Session>('/login/congressista', { method: 'POST', body: JSON.stringify({ email, senha }) })
}

export function registerCongressista(payload: { email: string; senha: string; nome: string; cpf: string; dataNascimento: string }) {
  return request<{ mensagem: string }>('/cadastro/congressista', { method: 'POST', body: JSON.stringify(payload) })
}

export function getProfile(session: Session) { return request<Profile>(`/usuarios/${session.id}`, {}, session) }
export function getEvents(session: Session, organizerId: number) { return request<EventSummary[]>(`/eventos/organizador/${organizerId}`, {}, session) }
export function getTicketTypes(session: Session) { return request<TicketType[]>('/eventos/tipos-ingresso', {}, session) }
export function getEvent(session: Session, id: number) { return request<EventDetails>(`/eventos/${id}`, {}, session) }
export function getAvailableEvents(session: Session) { return request<AvailableEventGroup[]>('/eventos/disponiveis', {}, session) }
export function getAvailableEvent(session: Session, id: number) { return request<EventDetails>(`/eventos/disponiveis/${id}`, {}, session) }
export function createEvent(session: Session, payload: EventPayload) { return request<EventDetails>('/eventos', { method: 'POST', body: JSON.stringify(payload) }, session) }
export function updateEvent(session: Session, id: number, payload: EventPayload) { return request<EventDetails>(`/eventos/${id}`, { method: 'PUT', body: JSON.stringify(payload) }, session) }
export function publishEvent(session: Session, id: number) { return request<EventDetails>(`/eventos/${id}/publicar`, { method: 'PATCH' }, session) }
export function closeEvent(session: Session, id: number) { return request<EventDetails>(`/eventos/${id}/encerrar`, { method: 'PATCH' }, session) }
export function deleteEvent(session: Session, id: number) { return request<void>(`/eventos/${id}`, { method: 'DELETE' }, session) }
export function getCart(session: Session) { return request<Cart>('/carrinho', {}, session) }
export function addCartItem(session: Session, eventoIngressoId: number, quantidade: number) { return request<Cart>('/carrinho/itens', { method: 'POST', body: JSON.stringify({ evento_ingresso_id: eventoIngressoId, quantidade }) }, session) }
export function updateCartItem(session: Session, itemId: number, quantidade: number) { return request<Cart>(`/carrinho/itens/${itemId}`, { method: 'PATCH', body: JSON.stringify({ quantidade }) }, session) }
export function removeCartItem(session: Session, itemId: number) { return request<void>(`/carrinho/itens/${itemId}`, { method: 'DELETE' }, session) }
export function clearCart(session: Session) { return request<void>('/carrinho', { method: 'DELETE' }, session) }
export function createOrder(session: Session, idempotencyKey: string) {
  return request<Order>('/pedidos', { method: 'POST', headers: { 'Idempotency-Key': idempotencyKey } }, session)
}
export function startPayment(session: Session, orderId: number, idempotencyKey: string) {
  return request<Order>(`/pedidos/${orderId}/pagamento`, { method: 'POST', headers: { 'Idempotency-Key': idempotencyKey } }, session)
}
export function getOrder(session: Session, orderId: number) { return request<Order>(`/pedidos/${orderId}`, {}, session) }
export function simulatePayment(session: Session, orderId: number, approved: boolean) {
  return request<Order>(`/pedidos/${orderId}/pagamento/simular/${approved ? 'aprovar' : 'rejeitar'}`, { method: 'POST' }, session)
}
export function getMyEvents(session: Session) { return request<MyEvent[]>('/meus-eventos', {}, session) }
export function getMyEvent(session: Session, eventId: number) { return request<MyEvent>(`/meus-eventos/${eventId}`, {}, session) }
