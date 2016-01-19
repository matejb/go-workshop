# go-workshop

Izrada aplikcija je podjeljenja u sljedeće faze:

1. Izraditi cli app koji iz json datoteke čita listu css datoteke te datoteke spaja u jednu css datoteku, primjer pokretanja

```bash
app.exe --list moja_lista.txt --out merged.css
```

2. Nadograditi app da može pratiti promjene nad listom css datoteku automatski te čim se one dese generira finalnu css datoteku, primjer pokretanja

```bash
app.exe --watch on --list moja_lista.txt --out merged.css
```

3. Nadograditi app da može servirati finalnu css datoteku na proizvoljnom portu, primjer pokretanja

```bash
app.exe --watch on --serve 8080 --list moja_lista.txt --out merged.css
```
