# Go Developer Competency Matrix

## Core Go Development

### Fundamentals
| Competency | Level 1 (Junior) | Level 2 (Middle) | Level 3 (Senior) | Level 4 (Expert) | Weight |
|------------|------------------|------------------|------------------|------------------|---------|
| **Go Syntax & Basics** | Basic syntax, variables, data types, operators, conditions, loops | Methods, interfaces, packages, modules, error handling | Reflection, testing, documentation, meta-programming | Performance optimization, profiling, mentoring | 1.0 |
| **Data Structures** | Arrays, slices, maps creation and iteration | Structs, embedding, custom data structures | Trees, graphs, caching, complex structures | Parallel data structures, optimization, mentoring | 1.1 |
| **Memory Management** | GC principles, pointers, null pointers | Memory optimization, object pools, profiling | GC tuning, lifecycle management, thread-safe memory | Runtime analysis, advanced tools, mentoring | 1.1 |
| **Concurrency** | Goroutines, channels, basic synchronization | Mutex, WaitGroup, patterns, context | Profiling, optimization, distributed systems | Runtime optimization, library development, mentoring | 1.3 |

### Advanced Go
| Competency | Level 1 | Level 2 | Level 3 | Level 4 | Weight |
|------------|---------|---------|---------|---------|---------|
| **HTTP & Web** | Basic HTTP server, routing, requests | Middlewares, JSON, forms, context | Optimization, security, testing | Microservices, API gateway, load balancing | 1.0 |
| **Quality Assurance** | Unit testing, coverage, basic docs | Mocking, integration tests, linting | CI/CD, profiling, code quality | Test automation, metrics, mentoring | 1.1 |
| **OS Interaction** | File I/O, command execution, paths | Environment, signals, permissions | Cross-platform, syscalls, performance | Optimization, system tools, mentoring | 1.0 |
| **gRPC** | Basic concepts, .proto files, simple implementation | Streaming, interceptors, testing | Security, performance, monitoring | Microservices, patterns, mentoring | 1.2 |

## System Design & Architecture

### Architecture Fundamentals
| Competency | Level 1 | Level 2 | Level 3 | Level 4 | Weight |
|------------|---------|---------|---------|---------|---------|
| **System Design** | Component architecture, monolith vs microservices | Synchronous/asynchronous communication | Scalable component architectures | Complex system architecture, technology trends | 1.3 |
| **Microservices** | Pros/cons understanding | Practical experience | Development with pros/cons consideration | Design and optimization of complex systems | 1.2 |
| **Containerization** | Basic orchestration concepts | kubectl, docker, scaling principles | High-load system management | Orchestrator architecture and principles | 1.1 |
| **Reliability** | Basic distributed transactions, scalability | Transaction optimization, fault tolerance | Various transaction patterns, auto-scaling | Advanced transactions, load balancing | 1.2 |

### Performance & Scalability
| Competency | Level 1 | Level 2 | Level 3 | Level 4 | Weight |
|------------|---------|---------|---------|---------|---------|
| **Performance** | Basic optimization, memory management | Profiler usage, application optimization | Database optimization, scalable design | Scalable architectures, cloud computing | 1.1 |
| **Latency & Throughput** | Basic concepts, impact understanding | Analysis and optimization | Network optimization, resource management | Fine-tuning, distributed systems | 1.1 |
| **Availability & Consistency** | Concurrency basics, simple consistency | Transactions, ACID, state consistency | Data consistency in distributed systems | Architectural patterns, failover design | 1.2 |

## Software Engineering Practices

### Development Methodologies
| Competency | Level 1 | Level 2 | Level 3 | Level 4 | Weight |
|------------|---------|---------|---------|---------|---------|
| **SDLC** | SDLC principles understanding | All stages understanding | Deep process understanding | Modern methodologies and technologies | 1.0 |
| **Requirements** | Basic requirements management, decomposition | Iterative development, documentation | Deep decomposition understanding | Effective decomposition in complex projects | 1.0 |
| **CI/CD** | Principles and benefits understanding | Setup and maintenance experience | Complex processes development | Automation and optimization in large projects | 1.1 |
| **Git & Version Control** | Basic concepts, theoretical knowledge | Basic techniques, GitFlow understanding | Advanced workflows, merge strategies | Team leadership, methodology standardization | 1.0 |

### Software Design
| Competency | Level 1 | Level 2 | Level 3 | Level 4 | Weight |
|------------|---------|---------|---------|---------|---------|
| **Coding Best Practices** | Maintainability, readability, formatting | Informative comments, descriptive naming | OOP principles, dependency minimization | Architectural patterns, performance optimization | 1.2 |
| **OOP Principles** | Basic OOP concepts, class-object differences | Deep OOP understanding, interface usage | Class design expertise, performance analysis | Advanced OOP mastery | 1.1 |
| **Design Principles** | DRY, KISS, SOLID description | SOLID application, DRY implementation | High-performance system creation | Architectural cleanliness, effective solutions | 1.2 |
| **Design Patterns** | Theoretical knowledge of basic patterns | Pattern differences, practical application | Multiple patterns, group understanding | New pattern creation | 1.3 |

## Security

### Web Security
| Competency | Level 1 | Level 2 | Level 3 | Level 4 | Weight |
|------------|---------|---------|---------|---------|---------|
| **Authentication & Authorization** | Basic concepts, protocols, sessions | Multi-factor auth, secure processes, RBAC | Data security, system integration | Secure mechanism design | 1.2 |
| **Security Controls** | Basic principles, SAST/DAST understanding | Results analysis, vulnerability identification | Strategy development based on scanning | Team consultation, standards development | 1.1 |
| **Data Security** | Basic threats, encryption methods | Encryption techniques, vulnerability resolution | Advanced methods, distributed systems | Organizational strategy development | 1.3 |
| **OWASP Risks** | OWASP Top 10 knowledge | Secure architecture development | Specific threat knowledge, risk assessment | Best practices implementation, team leadership | 1.2 |

## Question Type Examples

### Multiple Choice (Theory)
**Level: Middle | Competency: Concurrency**
```go
// What does this code output?
package main

func main() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    close(ch)
    
    for i := range ch {
        println(i)
    }
}
```
A) 1, 2
B) 1, 2, 0
C) Deadlock
D) 0, 0
## Coding Task
### Level: Senior | Competency: Concurrency
```go
// Implement a thread-safe cache with TTL
// Requirements:
// - Use sync.RWMutex
// - Goroutine for cleaning expired entries
// - Methods: Get, Set, Delete
type Cache struct {
    // Implement this
}

func NewCache() *Cache {
    // Your implementation
}
```
## Architecture Design
### Level: Expert | Competency: System Design
```text
Design a system for handling 1M requests per second with:
- Low latency requirements
- Fault tolerance
- Horizontal scalability
- Monitoring and logging
```
## Assessment Configuration Examples
### Backend Developer Role
```yaml
role: backend_developer
target_level: middle
competency_weights:
  go_fundamentals: 1.2
  concurrency: 1.4
  data_structures_go: 1.3
  system_design: 1.1
  databases: 1.2
question_distribution:
  multiple_choice: 8
  coding: 8
  architecture: 3
  debugging: 1
```
### Full-Stack Developer Role
```yaml
role: fullstack_developer
target_level: senior
competency_weights:
  go_fundamentals: 1.1
  http_go: 1.3
  concurrency: 1.2
  system_design: 1.4
  web_security: 1.3
```
