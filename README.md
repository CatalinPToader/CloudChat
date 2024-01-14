
# CloudChat

CloudChat este o aplicatie de mesagerie. Functionalitatile pe care utilizatorii le pot accesa sunt urmatoarele: autentificare
(utilizand adresa de mail, un nume de utilizator si parola), logare si conectare pe canale de chat specifice pentru a comunica prin mesaje.

## Echipa

 - Dragomir Constantin-Cristian
 - Toader Petru Catalin
 - Stan Sabina

## Structura microserviciilor CloudChat

Aplicatia este implementata pe un cluster Kubernetes, creat folosind utilitarul Kind. Acest cluster este configurat si gestionat prin Terraform si este format dintr-un control plane si 2 workeri. In cadrul clusterului, Portainer este instalat ca administrator, avand drepturi depline asupra clusterului prin asignarea tuturor permisiunilor in cadrul rolului sau in cluster. Pentru deployment-ul aplicatiilor s-au creat doua fisiere Terraform: primul fisier este echivalentul unui Deployment in Kubernetes, iar cel de-al doilea fisier gestioneaza microserviciile, facilitand expunerea porturilor containerelor in mod corespunzator.

Aplicatia este formata din:
 - un microserviciu de autentificare si autorizare: imagine de Docker cu Flask API si Firebase Authentication
 - un microserviciu de "business logic": imagine de Docker cu un REST API scris in Go, si un frontend in HTML cu Javascript
 - un microserviciu de tip baza de date: PostgreSQL
 - un microserviciu de gestiune a bazelor de date: pgAdmin
 - un microserviciu de tip utilitar grafic de gestiune a clusterului: Portainer

Pentru a incorpora "business logic" (backend) si serviciul de autentificare s-au realizat imagini de Docker. Pentru a putea proviziona nodurile cu aceste imagini a fost necesara crearea unei imagini de Docker separata de cluster care functioneaza ca un registru privat de imagini, similar cu Dockerhub. Acest registru a fost inregistrat ca o sursa in cadrul clusterului.
De asemenea, pentru pgAdmin si PostgreSQL s-au construit diverse fisiere aditionale pentru configurarea unui server si crearea de tabele in baza de date.

## Cerinte preliminare

Este necesar sa aveti instalate pe propriul dispozitiv urmatoarele:

 - Docker
 - Kind - pentru crearea clusterului Kubernetes
 - Terraform - pentru aplicarea configuratiei pe cluster

## Rulare
Toate aspectele legate de implementare sunt gestionate printr-un script numit build-infra.sh care faciliteaza procesul de configurare a intregii infrastructuri, incluzand Kubernetes, registrul privat de Docker, crearea imaginilor pentru backend si autentificare, precum si aplicarea configuratiei Terraform, toate acestea fiind realizate printr-o singura comanda.

**_Note:_** Este posibil sa fie necesara rularea scripturilor cu **sudo**

Din folderul radacina al proiectului se executa comanda:
```bash
  ./build-infra.sh
```

## Oprire deployment si eliberare resurse

Din folderul radacina al proiectului se executa comanda:
```bash
  ./destroy-infra.sh
```
