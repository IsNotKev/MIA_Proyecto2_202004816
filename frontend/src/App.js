
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

  function verRep(){
    var obj = { 'cmd': document.getElementById('Instrucciones').value }
      fetch(`http://localhost:5000/reporte`, {
        method: 'POST',
        body: JSON.stringify(obj),
      })
        .then(res => res.json())
        .catch(err => {
          console.error('Error:', err)
          alert("Ocurrio un error, ver la consola")
        })
        .then(response => {
          var image = new Image();
          image.src = response.result
          document.getElementById('RepDiv').appendChild(image);
        })
  }

  return (
    <div className="App" style={{width:'80%', margin:'auto', marginTop:'2%'}}>
      <nav class="navbar navbar-expand-lg navbar-dark bg-dark" style={{ width: '100%' }}>
          <div class="container-fluid">
            <a class="navbar-brand" href="#">Proy2</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
              <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarSupportedContent">
              <ul class="navbar-nav me-auto mb-2 mb-lg-0">
              <li class="nav-item">
                  <a class="nav-link" style={{ cursor: 'pointer' }} onClick={ejecutar}>Ejecutar</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" style={{ cursor: 'pointer' }}data-bs-toggle="modal" data-bs-target="#modalLogin">Login</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" style={{ cursor: 'pointer' }} data-bs-toggle="modal" data-bs-target="#modalRep" onClick={verRep}>Reportes</a>
                </li>
              </ul>
            </div>
          </div>
        </nav>
      <input required class="form-control" type="file" id="file-input" onChange={leerArchivo} style={{ marginTop: '2%', width: '50%' }}></input>
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


      <div class="modal fade" id="modalRep" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog modal-fullscreen">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="exampleModalLabel">REPORTES</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
              <div id="RepDiv"></div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            </div>
          </div>
        </div>
      </div>

      <div class="modal fade" id="modalLogin" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog modal-fullscreen">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="exampleModalLabel">LOGIN</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body" style={{width:'50%', margin:'auto'}}>
              <div class="mb-3 row">
                <label for="inputPassword" class="col-sm-2 col-form-label">Disco: </label>
                <div class="col-sm-10">
                  <input type="text" class="form-control" id="inputDisco"></input>
                </div>
              </div>
              <div class="mb-3 row">
                <label for="inputPassword" class="col-sm-2 col-form-label">Usuario:</label>
                <div class="col-sm-10">
                  <input type="text" class="form-control" id="inputUser"></input>
                </div>
              </div>
              <div class="mb-3 row">
                <label for="inputPassword" class="col-sm-2 col-form-label">Contraseña:</label>
                <div class="col-sm-10">
                  <input type="text" class="form-control" id="inputContra"></input>
                </div>
              </div>
              <div class="col-auto">
                <button type="submit" class="btn btn-primary mb-3">Ingresar</button>
              </div>

            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            </div>
          </div>
        </div>
      </div>

    </div>
  );
}

export default App;
