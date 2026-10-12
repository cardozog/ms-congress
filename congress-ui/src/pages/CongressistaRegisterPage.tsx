import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { registerCongressista } from '../services/api'
import { errorMessage } from '../utils'

export function CongressistaRegisterPage() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ nome: '', email: '', senha: '', cpf: '', dataNascimento: '' })
  const [feedback, setFeedback] = useState('')
  const [loading, setLoading] = useState(false)
  function update(field: keyof typeof form, value: string) { setForm(current => ({ ...current, [field]: value })) }
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setFeedback(''); setLoading(true)
    try { await registerCongressista(form); navigate('/congressista/login', { replace: true, state: { registered: true } }) }
    catch (error) { setFeedback(errorMessage(error, 'Não foi possível concluir o cadastro.')) }
    finally { setLoading(false) }
  }
  return <><header className="topbar public-topbar"><a className="brand" href="/congressista/login"><span className="brand-mark">MS</span><span><strong>MS Congress</strong><small>área do congressista</small></span></a></header><section className="form-view auth-view"><div className="form-view-heading"><span className="kicker">Novo congressista</span><h2>Crie sua conta</h2><p className="section-description">Preencha seus dados para acessar a área de eventos.</p></div><form className="form-panel register-panel" onSubmit={submit}><label htmlFor="register-name">Nome completo</label><input id="register-name" value={form.nome} onChange={(event) => update('nome', event.target.value)} required /><label htmlFor="register-email">E-mail</label><input id="register-email" type="email" value={form.email} onChange={(event) => update('email', event.target.value)} required /><label htmlFor="register-cpf">CPF</label><input id="register-cpf" value={form.cpf} onChange={(event) => update('cpf', event.target.value)} required /><label htmlFor="register-birth">Data de nascimento</label><input id="register-birth" type="date" value={form.dataNascimento} onChange={(event) => update('dataNascimento', event.target.value)} /><label htmlFor="register-password">Senha</label><input id="register-password" type="password" minLength={8} value={form.senha} onChange={(event) => update('senha', event.target.value)} required /><button className="primary-button full-width" disabled={loading}>{loading ? 'Criando conta...' : 'Criar conta'} <b>→</b></button>{feedback && <p className="feedback">{feedback}</p>}<p className="auth-link">Já possui conta? <Link to="/congressista/login">Entrar</Link></p></form></section></>
}
