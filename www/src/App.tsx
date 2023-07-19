import { useQuery } from "@apollo/client";
import { gql } from "./__generated__/gql";
import { useEffect, useState } from "react";

const GET_SETTINGS = gql(`
  query RootQuery {
    settings {
      LibreLinkUpUsername
      LibreLinkUpEndpoint
    }
  }
`);

function App() {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const { loading, error, data } = useQuery(GET_SETTINGS);

  useEffect(() => {
    if (!loading && data) {
      setUsername(data.settings.LibreLinkUpUsername);
    }
  }, [loading, data]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error : {error.message}</p>;

  function onChange(username: string, password: string, endpoint: string) {
    console.log("onChange: ", username, password, endpoint);
    setUsername(username);
    setPassword(password);
  }

  return (
    <>
      {
        data &&
        <div className="">
          <LibreLinkupForm
            username={username}
            password={password}
            endpoint={username}
            onChange={onChange}
          />
        </div>
      }
    </>
  )
}

interface LibreLinkupFormProps {
  username: string
  password: string
  endpoint: string
  onChange: (username: string, password: string, endpoint: string) => void
}

function LibreLinkupForm(props: LibreLinkupFormProps) {
  return (
    <div className="p-6 max-w-sm mx-auto bg-white rounded-xl shadow-lg items-center space-x-4">
      <div className="shrink-0 text-xl font-medium mb-4">
        LibreLinkUp Settings
      </div>
      <div>
        <label className="font-medium text-slate-700" htmlFor="username">Username</label>
      </div>
      <div>
        <input
          className="border-solid border-2 border-sky-500 rounded"
          name="username"
          type="text"
          value={props.username}
          onChange={(e) => props.onChange(e.target.value, "", props.endpoint)}
          placeholder="Username"
        />
      </div>
      <div>
        <label className="font-medium text-slate-700" htmlFor="password">Password</label>
      </div>
      <div>
        <input
          className="border-solid border-2 border-sky-500 rounded"
          name="password"
          type="password"
          value={props.password}
          onChange={(e) => props.onChange(props.username, e.target.value, props.endpoint)}
          placeholder="Password"
        />
      </div>
      <div>
        <label className="font-medium text-slate-700" htmlFor="endpoint">Endpoint</label>
      </div>
      <div>
        <input
          className="border-solid border-2 border-sky-500 rounded"
          name="endpoint"
          type="text"
          value={props.endpoint}
          onChange={(e) => props.onChange(props.username, props.password, e.target.value)}
          placeholder="Endpoint"
        />
      </div>
    </div>
  )
}

export default App
