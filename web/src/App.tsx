function App() {
  return (
    <>
      <Header title="Freelance Marketplace" />
    </>
  );
}

function Header(props: { title: string }) {
  return <h1>{props.title}</h1>;
}

export default App;
