import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { createOrder, getCart, startPayment } from '../services/api'
import { useSession } from '../session/useSession'
import type { Cart, Order } from '../types'
import { errorMessage } from '../utils'

const currency = (value: number) => value.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })

export function CongressistaCheckoutPage() {
  const { session } = useSession()
  const navigate = useNavigate()
  const [cart, setCart] = useState<Cart | null>(null)
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [feedback, setFeedback] = useState('')
  const [startingPayment, setStartingPayment] = useState(false)
  const idempotencyKey = useMemo(() => {
    const storageKey = `ms-congress-checkout-key-${session?.id ?? 'anonymous'}`
    const key = sessionStorage.getItem(storageKey)
    if (key) return key
    const next = crypto.randomUUID()
    sessionStorage.setItem(storageKey, next)
    return next
  }, [session?.id])

  useEffect(() => {
    if (!session) return
    void getCart(session).then(setCart).catch(error => setFeedback(errorMessage(error, 'Não foi possível carregar o resumo da compra.'))).finally(() => setLoading(false))
  }, [session])

  async function finishOrder() {
    if (!session) return
    setSubmitting(true)
    setFeedback('')
    try {
      const created = await createOrder(session, idempotencyKey)
      navigate(`/congressista/pagamento/${created.id}`)
    } catch (error) {
      setFeedback(errorMessage(error, 'Não foi possível criar o pedido. Revise a disponibilidade dos ingressos.'))
    } finally {
      setSubmitting(false)
    }

  }

  async function beginPayment() {
    if (!session || !order) return
    setStartingPayment(true)
    try { setOrder(await startPayment(session, order.id, idempotencyKey)) }
    catch (error) { setFeedback(errorMessage(error, 'Não foi possível iniciar o pagamento.')) }
    finally { setStartingPayment(false) }
  }

  if (loading) return <div className="app-shell congressista-shell"><main className="main-content"><div className="empty-state"><p>Carregando checkout...</p></div></main></div>
  if (order) return <div className="app-shell congressista-shell"><main className="main-content"><section className="checkout-view"><span className="kicker">Pedido criado</span><h1>Pedido #{order.id}</h1><p>Seu pedido está aguardando confirmação de pagamento. Os ingressos ainda não foram adquiridos.</p>{feedback && <p className="feedback">{feedback}</p>}<div className="cart-total"><span>Status: {order.status}</span><strong>{currency(order.total)}</strong></div>{order.pagamento?.referenciaExterna ? <p className="feedback">Pagamento iniciado. Aguardando confirmação segura do provedor.</p> : <button className="primary-button" type="button" disabled={startingPayment} onClick={() => void beginPayment()}>{startingPayment ? 'Iniciando pagamento...' : 'Iniciar pagamento'}</button>}<button className="secondary-button" type="button" onClick={() => navigate('/congressista/eventos')}>Voltar aos eventos</button></section></main></div>
  if (!cart?.itens.length) return <div className="app-shell congressista-shell"><main className="main-content"><section className="empty-state"><h3>Seu carrinho está vazio</h3><Link className="primary-button" to="/congressista/eventos">Descobrir eventos</Link></section></main></div>

  return <div className="app-shell congressista-shell"><main className="main-content"><header className="topbar"><div><p className="eyebrow">Área do congressista</p><h1>Checkout</h1></div><div className="user-area"><strong>{session?.email}</strong><small>Congressista</small></div></header><section className="checkout-view"><div className="checkout-grid"><div><span className="kicker">Resumo da compra</span>{feedback && <p className="feedback">{feedback}</p>}{cart.itens.map(item => <article className="cart-item" key={item.id}><div><span className="kicker">{item.evento.nome}</span><h3>{item.ingresso.nome}</h3><p>{item.quantidade} × {currency(item.precoUnitario)}</p></div><strong>{currency(item.subtotal)}</strong></article>)}</div><aside className="checkout-summary"><span className="kicker">Congressista</span><p>{session?.email}</p><div className="cart-total"><span>Total</span><strong>{currency(cart.valorTotal)}</strong></div><button className="primary-button" type="button" disabled={submitting} onClick={() => void finishOrder()}>{submitting ? 'Criando pedido...' : 'Finalizar pedido'}</button><small>O pedido ficará pendente até a confirmação do pagamento.</small></aside></div></section></main></div>
}
