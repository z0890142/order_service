import './App.css';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Login from './components/Login';
import PatientList from './components/PatientList';
import PrivateRoute from './components/PrivateRoute';

function App() {
  return (
    <div className="App">
        <Routes>
          <Route element={<Login />} path={'/'}></Route>
          <Route path="/patients" element={<PrivateRoute Component={PatientList} />} />

        </Routes>
     
   </div>
  );
}

export default App;
