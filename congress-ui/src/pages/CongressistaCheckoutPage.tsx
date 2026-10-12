import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { createOrder, getCart, startPayment } from '../services/api'
import { useSession } from '../session/useSession'
import type { Cart } from '../types'
import { errorMessage } from '../utils'

const currency = (value: number) => value.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })

export function CongressistaCheckoutPage() {
  const { session } = useSession()
  const navigate = useNavigate()
  const [cart, setCart] = useState<Cart | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [feedback, setFeedback] = useState('')
  const [card, setCard] = useState({ numeroCartao: '', cpf: '', vencimento: '', codigoSeguranca: '' })
  const [idempotencyKey] = useState(() => crypto.randomUUID())
  const [createdOrderId, setCreatedOrderId] = useState<number | null>(null)

  useEffect(() => {
    if (!session) return
    void getCart(session).then(setCart).catch(error => setFeedback(errorMessage(error, 'Não foi possível carregar o resumo da compra.'))).finally(() => setLoading(false))
  }, [session])

  async function finishOrder() {
    if (!session) return
    setSubmitting(true)
    setFeedback('')
    try {
      const orderId = createdOrderId ?? (await createOrder(session, idempotencyKey)).id
      setCreatedOrderId(orderId)
      await startPayment(session, orderId, idempotencyKey, card)
      navigate(`/congressista/pagamento/${orderId}`)
    } catch (error) {
      setFeedback(errorMessage(error, 'Não foi possível criar o pedido. Revise a disponibilidade dos ingressos.'))
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <div className="app-shell congressista-shell"><main className="main-content"><div className="empty-state"><p>Carregando checkout...</p></div></main></div>
  if (!cart?.itens.length) return <div className="app-shell congressista-shell"><main className="main-content"><section className="empty-state"><h3>Seu carrinho está vazio</h3><Link className="primary-button" to="/congressista/eventos">Descobrir eventos</Link></section></main></div>

  const isCardComplete = Object.values(card).every(value => value.trim() !== '')

  return <div className="app-shell congressista-shell"><main className="main-content"><header className="topbar"><div><p className="eyebrow">Área do congressista</p><h1>Checkout</h1></div><div className="user-area"><strong>{session?.email}</strong><small>Congressista</small></div></header><section className="checkout-view"><div className="checkout-grid"><div><span className="kicker">Resumo da compra</span>{feedback && <p className="feedback">{feedback}</p>}{cart.itens.map(item => <article className="cart-item" key={item.id}><div><span className="kicker">{item.evento.nome}</span><h3>{item.ingresso.nome}</h3><p>{item.quantidade} × {currency(item.precoUnitario)}</p></div><strong>{currency(item.subtotal)}</strong></article>)}<form className="form-panel payment-form" onSubmit={event => { event.preventDefault(); if (isCardComplete) void finishOrder() }}><span className="kicker">Dados do cartão</span><p className="section-description">O número completo e o código de segurança não são armazenados.</p><div className="form-grid"><div className="field-wide"><label htmlFor="checkout-card-number">Número do cartão</label><input id="checkout-card-number" inputMode="numeric" autoComplete="cc-number" value={card.numeroCartao} onChange={event => setCard({ ...card, numeroCartao: event.target.value })} required /></div><div><label htmlFor="checkout-card-cpf">CPF</label><input id="checkout-card-cpf" inputMode="numeric" autoComplete="off" value={card.cpf} onChange={event => setCard({ ...card, cpf: event.target.value })} required /></div><div><label htmlFor="checkout-card-expiry">Vencimento</label><input id="checkout-card-expiry" placeholder="MM/AA" inputMode="numeric" autoComplete="cc-exp" value={card.vencimento} onChange={event => setCard({ ...card, vencimento: event.target.value })} required /></div><div><label htmlFor="checkout-card-cvv">Código de segurança</label><input id="checkout-card-cvv" type="password" inputMode="numeric" autoComplete="cc-csc" maxLength={4} value={card.codigoSeguranca} onChange={event => setCard({ ...card, codigoSeguranca: event.target.value })} required /></div></div><button className="primary-button full-width" type="submit" disabled={submitting || !isCardComplete}>{submitting ? 'Enviando pagamento...' : 'Finalizar pedido'}</button></form></div><aside className="checkout-summary"><span className="kicker">Congressista</span><p>{session?.email}</p><div className="cart-total"><span>Total</span><strong>{currency(cart.valorTotal)}</strong></div><small>Após o envio, você poderá aprovar ou reprovar o pagamento na tela de demonstração.</small></aside></div></section></main></div>
}
