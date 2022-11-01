
import './App.css';

function App() {

  function leerArchivo(e) {
    var archivo = e.target.files[0];
    if (!archivo) {
      return;
    }
    var lector = new FileReader();
    lector.onload = function (e) {
      var contenido = e.target.result;
      mostrarContenido(contenido);
    };
    lector.readAsText(archivo);
  }

  function mostrarContenido(contenido) {
    document.getElementById('Instrucciones').value = contenido;
  }

  function ejecutar(){
    if (document.getElementById('Instrucciones').value !== ""){
      var obj = { 'cmd': document.getElementById('Instrucciones').value }
      fetch(`http://localhost:5000/ejecutar`, {
        method: 'POST',
        body: JSON.stringify(obj),
      })
        .then(res => res.json())
        .catch(err => {
          console.error('Error:', err)
          alert("Ocurrio un error, ver la consola")
        })
        .then(response => {
          document.getElementById('Resultado').value = response.result;
        })
    }else{
      alert('Consola Vacía')
    }
    
  }

  return (
    <div className="App" style={{width:'80%', margin:'auto', marginTop:'2%'}}>

      <h1 style={{marginBottom:'10%'}}> Cargar Instrucciones </h1>
      <input required class="form-control" type="file" id="file-input" onChange={leerArchivo} style={{ marginTop: '2%', width: '50%' }}></input>
      <br></br>
      <button type="button" class="btn btn-success" style={{marginLeft:'-92%'}} onClick={ejecutar}>Ejecutar</button>
      <br></br>
      <div class="form-floating" style={{marginTop:'2%'}}>
        <textarea class="form-control" id="Instrucciones" style={{height:'200px'}}></textarea>
        <label for="floatingTextarea">Instrucciones</label>
      </div>
      <br></br>
      <div class="form-floating" style={{marginTop:'2%'}}>
        <textarea class="form-control" id="Resultado" style={{height:'200px'}} disabled></textarea>
        <label for="floatingTextarea">Resultado</label>
      </div>
    </div>
  );
}

export default App;
