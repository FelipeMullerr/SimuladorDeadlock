package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RecursoCompartilhado struct {
	nome string
	dono *Transacao
}
type Transacao struct {
	id          int
	timestamp   int64
	foiAbortada bool
}

type EsperaPorRecurso struct {
	Transacao *Transacao
	Recurso   *RecursoCompartilhado
}

var (
	recursoX      = &RecursoCompartilhado{nome: "X"}
	recursoY      = &RecursoCompartilhado{nome: "Y"}
	grupoDeEspera sync.WaitGroup
	mutexGlobal   sync.Mutex
	listaDeEspera []EsperaPorRecurso
	transacoesFinalizadas = make(map[int]bool)
)

func printColor(cor string, mensagem string, args ...interface{}) {
	var colorCode string
	switch cor {
	case "red":
		colorCode = "\033[31m"
	case "yellow":
		colorCode = "\033[33m"
	case "green":
		colorCode = "\033[32m"
	case "blue":
		colorCode = "\033[38;5;19m"
	case "magenta":
		colorCode = "\033[35m"
	case "cyan":
		colorCode = "\033[36m"
	case "white":
		colorCode = "\033[37m"
	case "boldRed":
		colorCode = "\033[1;31m"
	case "underlineBlue":
		colorCode = "\033[4;34m"
	case "darkBlue":
		colorCode = "\033[38;5;19m"
	default:
		colorCode = "\033[0m"
	}
	fmt.Printf(colorCode+mensagem+"\033[0m\n", args...)
}

func executarTransacao(transacao *Transacao) {
	defer grupoDeEspera.Done()

	printColor("reset", "Transação [%d] iniciou sua execução.", transacao.id)
	rand_t()

	if verificarAbortada(transacao) {
		return
	}

	if transacao.id%2 == 0 {
		// Ordem: X depois Y
		if !lock_item(transacao, recursoX) {
			return
		}

		if verificarAbortada(transacao) {
			return
		}

		rand_t()

		if !lock_item(transacao, recursoY) {
			return
		}
	} else {
		if verificarAbortada(transacao) {
			return
		}
		// Ordem: Y depois X
		if !lock_item(transacao, recursoY) {
			return
		}

		rand_t()

		if !lock_item(transacao, recursoX) {
			return
		}
	}

	rand_t()

	if verificarAbortada(transacao) {
		return
	}

	unlock_item(transacao, recursoX)
	rand_t()
	unlock_item(transacao, recursoY)
	rand_t()

	mutexGlobal.Lock()
	if !transacoesFinalizadas[transacao.id] {
		transacoesFinalizadas[transacao.id] = true
		printColor("underlineBlue", "Transação [%d] finalizou (commit) com sucesso.", transacao.id)
	}
	mutexGlobal.Unlock()
}

func detectorDeDeadlock() {
	for {
		time.Sleep(500 * time.Millisecond)

		mutexGlobal.Lock()
		if len(listaDeEspera) >= 2 {
			for i := 0; i < len(listaDeEspera); i++ {
				for j := i + 1; j < len(listaDeEspera); j++ {
					t1 := listaDeEspera[i].Transacao
					r1 := listaDeEspera[i].Recurso
					t2 := listaDeEspera[j].Transacao
					r2 := listaDeEspera[j].Recurso

					if r1.dono == t2 && r2.dono == t1 {
						var abortar *Transacao
						// Aborta a transação com o maior timestamp (mais antiga)
						if t1.timestamp > t2.timestamp {
							abortar = t1
						} else {
							abortar = t2
						}
						printColor("red", "Deadlock detectado entre transações [%d] e [%d].", t1.id, t2.id)
						printColor("red", "Abortando transação [%d].", abortar.id)
						abortar.foiAbortada = true
						liberarRecursos(abortar)
						go reiniciarTransacao(abortar)
					}
				}
			}
		}
		mutexGlobal.Unlock()
	}
}

func liberarRecursos(t *Transacao) {
	for _, r := range []*RecursoCompartilhado{recursoX, recursoY} {
		if r.dono == t {
			r.dono = nil
			printColor("blue", "Transação [%d] liberou recurso %s.", t.id, r.nome)
		}
	}
	novaLista := []EsperaPorRecurso{}
	for _, e := range listaDeEspera {
		if e.Transacao != t {
			novaLista = append(novaLista, e)
		}
	}
	listaDeEspera = novaLista
}

func verificarAbortada(t *Transacao) bool {
	if t.foiAbortada {
		printColor("red", "Transação [%d] abortada no meio da execução.", t.id)
		return true
	}
	return false
}

func reiniciarTransacao(t *Transacao) {
	mutexGlobal.Lock()
	if transacoesFinalizadas[t.id] {
		mutexGlobal.Unlock()
		return
	}
	mutexGlobal.Unlock()

	novaTransacao := &Transacao{
		id:        t.id,
		timestamp: time.Now().UnixNano(),
	}
	printColor("blue", "Reiniciando transação [%d].", novaTransacao.id)
	grupoDeEspera.Add(1)
	go executarTransacao(novaTransacao)
}

func removerDaListaDeEspera(t *Transacao, r *RecursoCompartilhado) {
	novaLista := []EsperaPorRecurso{}
	for _, espera := range listaDeEspera {
		if espera.Transacao != t || espera.Recurso != r {
			novaLista = append(novaLista, espera)
		}
	}
	listaDeEspera = novaLista
}

func lock_item(transacao *Transacao, recurso *RecursoCompartilhado) bool {
	mutexGlobal.Lock()
	if recurso.dono == nil {
		recurso.dono = transacao
		mutexGlobal.Unlock()
		printColor("green", "Transação [%d] obteve o recurso %s.", transacao.id, recurso.nome)
		return true
	}
	// recurso ocupado, transação espera até o recurso estiver livre
	listaDeEspera = append(listaDeEspera, EsperaPorRecurso{transacao, recurso})
	printColor("yellow", "Transação [%d] está esperando recurso %s.", transacao.id, recurso.nome)
	mutexGlobal.Unlock()

	for {
		time.Sleep(50 * time.Millisecond)

		// verifica se a transação foi abortada
		mutexGlobal.Lock()
		if transacao.foiAbortada {
			mutexGlobal.Unlock()
			return false
		}
		if recurso.dono == nil {
			recurso.dono = transacao
			removerDaListaDeEspera(transacao, recurso)
			mutexGlobal.Unlock()
			printColor("green", "Transação [%d] obteve recurso %s após esperar.", transacao.id, recurso.nome)
			return true
		}
		mutexGlobal.Unlock()
	}
}

func unlock_item(transacao *Transacao, recurso *RecursoCompartilhado) {
	mutexGlobal.Lock()
	if recurso.dono == transacao {
		recurso.dono = nil
		printColor("blue", "Transação [%d] liberou o recurso %s.", transacao.id, recurso.nome)
	}
	mutexGlobal.Unlock()
}

func rand_t() {
	tempo := rand.Intn(300)
	time.Sleep(time.Millisecond * time.Duration(tempo))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go detectorDeDeadlock()

	// Criacao de cada thread que ira executar a funcao "executarTransacao"
	for i := 0; i < 5; i++ {
		transacao := &Transacao{id: i, timestamp: time.Now().UnixNano()}
		grupoDeEspera.Add(1)
		go executarTransacao(transacao)
		// Aguarda um tempo aleatório antes de inicializar a próxima transação
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}

	grupoDeEspera.Wait()
}
