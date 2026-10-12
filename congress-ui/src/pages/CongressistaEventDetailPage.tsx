import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { addCartItem, getAvailableEvent } from '../services/api'
import { useSession } from '../session/useSession'
import type { EventDetails } from '../types'
import { errorMessage, formatDate } from '../utils'

export function CongressistaEventDetailPage() {
  const { id } = useParams()
  const { session } = useSession()
  const navigate = useNavigate()
  const [event, setEvent] = useState<EventDetails | null>(null)
  const [feedback, setFeedback] = useState('')
  const [quantities, setQuantities] = useState<Record<number, number>>({})
  const [added, setAdded] = useState(false)
  const [adding, setAdding] = useState(false)
  const [showAddedModal, setShowAddedModal] = useState(false)
  const [addErrorModal, setAddErrorModal] = useState('')
  async function refreshEvent() {
    if (!session || !id) return
    try {
      setEvent(await getAvailableEvent(session, Number(id)))
    } catch (error) {
      setFeedback(errorMessage(error, 'Não foi possível atualizar a disponibilidade dos ingressos.'))
    }
  }
  useEffect(() => {
    if (!session || !id) return
    const currentSession = session
    async function load() {
      try { setEvent(await getAvailableEvent(currentSession, Number(id))) }
      catch (error) { setFeedback(errorMessage(error, 'Não foi possível carregar o evento.')) }
    }
    void load()
  }, [session, id])
  const now = new Date()
  function isSaleOpen(ticket: EventDetails['ingressos'][number]) {
    return new Date(ticket.dataInicioVenda) <= now && new Date(ticket.dataFimVenda) >= now && ticket.disponivel > 0
  }
  function updateQuantity(ticketId: number, value: string, available: number) {
    const quantity = Math.max(0, Math.min(available, Number(value) || 0))
    setQuantities(current => ({ ...current, [ticketId]: quantity }))
    setAdded(false)
  }
  function addToCart() {
    if (!event) return
    const selected = event.ingressos.filter(ticket => (quantities[ticket.id] ?? 0) > 0)
    if (selected.length === 0) {
      setAddErrorModal('Selecione ao menos um ingresso para continuar.')
      setAdded(false)
      return
    }
    if (!session) return
    setAdding(true); setFeedback('')
    Promise.all(selected.map(ticket => addCartItem(session, ticket.id, quantities[ticket.id]))).then(() => {
      setAdded(true)
      setQuantities({})
      setShowAddedModal(true)
    }).catch(error => {
      setAdded(false)
      setAddErrorModal(errorMessage(error, 'Não foi possível adicionar os ingressos ao carrinho.'))
      void refreshEvent()
    }).finally(() => setAdding(false))
  }
  return <><section className="public-detail-view"><button className="back-button" type="button" onClick={() => navigate('/congressista/eventos')}>← Voltar para eventos</button>{event ? <><div className="public-event-cover">{event.eventoLogo ? <img src={event.eventoLogo} alt="" /> : <span>MS Congress</span>}</div><span className="kicker">Evento de {event.organizador.nomeFantasia || event.organizador.razaoSocial}</span><h2>{event.nome}</h2><p className="public-event-description">{event.descricao}</p><div className="public-event-info"><span><small>Começa em</small><strong>{formatDate(event.dataInicio, true)}</strong></span><span><small>Termina em</small><strong>{formatDate(event.dataFim, true)}</strong></span><span><small>Local</small><strong>{event.cidade} / {event.estado}</strong></span></div><p className="public-event-address">{event.logradouro}, {event.cep}</p><section className="public-tickets"><div className="public-tickets-heading"><div><span className="kicker">Ingressos</span><h3>Escolha como participar</h3></div><p>Os valores e a disponibilidade são atualizados pela organização.</p></div><div className="public-ticket-list">{event.ingressos.map(ticket => { const saleOpen = isSaleOpen(ticket); const quantity = quantities[ticket.id] ?? 0; return <div className={`public-ticket${saleOpen ? '' : ' is-unavailable'}`} key={ticket.id}><div><strong>{ticket.tipoIngresso.nome}</strong><span>{new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(ticket.preco)}</span><small>{saleOpen ? `${ticket.disponivel} disponível(is)` : ticket.disponivel === 0 ? 'Esgotado' : 'Vendas encerradas ou ainda não iniciadas'}</small></div><input type="number" min={0} max={ticket.disponivel} value={quantity} disabled={!saleOpen} onChange={inputEvent => updateQuantity(ticket.id, inputEvent.target.value, ticket.disponivel)} aria-label={`Quantidade de ${ticket.tipoIngresso.nome}`} /></div> })}</div><button className="primary-button add-cart-button" type="button" onClick={addToCart} disabled={adding}>{adding ? 'Adicionando...' : 'Adicionar ao carrinho'}</button>{feedback && <p className="feedback">{feedback}</p>}{added && <p className="success-feedback">Ingressos adicionados ao carrinho.</p>}</section><Link className="secondary-button" to="/congressista/eventos">Voltar à descoberta</Link></> : <div className="empty-state"><p>{feedback || 'Carregando detalhes...'}</p></div>}</section>{showAddedModal && <div className="flow-modal-backdrop"><section className="flow-modal" role="dialog" aria-modal="true" aria-labelledby="added-title"><span className="modal-icon">✓</span><h2 id="added-title">Ingressos adicionados</h2><p>Os ingressos de <strong>{event?.nome}</strong> foram adicionados ao seu carrinho.</p><button className="primary-button" type="button" onClick={() => navigate('/congressista/eventos', { replace: true })}>Voltar para eventos</button></section></div>}{addErrorModal && <div className="flow-modal-backdrop" onClick={() => setAddErrorModal('')}><section className="flow-modal" role="dialog" aria-modal="true" aria-labelledby="add-error-title" onClick={modalEvent => modalEvent.stopPropagation()}><span className="modal-icon is-error">!</span><h2 id="add-error-title">Não foi possível adicionar</h2><p>{addErrorModal}</p><button className="primary-button" type="button" onClick={() => setAddErrorModal('')}>Entendi</button></section></div>}</>
}
