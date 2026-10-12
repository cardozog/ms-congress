import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getCart, removeCartItem, updateCartItem } from '../services/api'
import { useSession } from '../session/useSession'
import type { Cart } from '../types'
import { errorMessage } from '../utils'

export function CongressistaCartPage() {
  const { session } = useSession(); const navigate = useNavigate()
  const [cart, setCart] = useState<Cart | null>(null); const [feedback, setFeedback] = useState('')
  useEffect(() => { if (!session) return; void getCart(session).then(setCart).catch(error => setFeedback(errorMessage(error, 'Não foi possível carregar o carrinho.'))) }, [session])
  async function update(id: number, quantity: number) { if (!session) return; try { setCart(await updateCartItem(session, id, quantity)) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível alterar a quantidade.')) } }
  async function remove(id: number) { if (!session) return; try { await removeCartItem(session, id); setCart(await getCart(session)) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível remover o item.')) } }
  return <div className="app-shell congressista-shell"><main className="main-content"><header className="public-nav"><Link className="public-brand" to="/congressista/eventos"><span className="brand-mark">MS</span><strong>MS Congress</strong></Link><nav><Link to="/congressista/eventos">Explorar eventos</Link><Link to="/meus-eventos">Meus Eventos</Link></nav><div className="public-actions"><span className="public-user">{session?.email}</span></div></header><section className="cart-view">{feedback && <p className="feedback">{feedback}</p>}{cart?.itens.length ? <><div className="cart-list">{cart.itens.map(item => <article className="cart-item" key={item.id}><div><span className="kicker">{item.evento.nome}</span><h3>{item.ingresso.nome}</h3><p>{item.precoUnitario.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })} por ingresso</p></div><input type="number" min={1} value={item.quantidade} onChange={event => void update(item.id, Number(event.target.value))} aria-label={`Quantidade de ${item.ingresso.nome}`} /><strong>{item.subtotal.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}</strong><button className="remove-ticket" type="button" onClick={() => void remove(item.id)} aria-label={`Remover ${item.ingresso.nome}`}>×</button></article>)}</div>  <div className="cart-total"><span>{cart.totalItens} item(ns)</span><strong>{cart.valorTotal.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}</strong></div><button className="primary-button" type="button" onClick={() => navigate('/congressista/checkout')}>Ir para checkout</button></> : <div className="empty-state"><h3>Seu carrinho está vazio</h3><p>Escolha ingressos em um evento para começar.</p><Link className="primary-button" to="/congressista/eventos">Descobrir eventos</Link></div>}</section></main></div>
}
