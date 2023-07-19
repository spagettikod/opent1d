import { useMutation, useQuery } from "@apollo/client";
import { gql } from "./__generated__/gql";
import { useEffect, useState } from "react";

const GET_SETTINGS = gql(`
  query GetSettings {
    settings {
      LibreLinkUpUsername
      LibreLinkUpPassword
    }
  }
`);

const SAVE_SETTINGS = gql(`
  mutation SaveSettings($username: String!, $password: String!) {
    saveSettings(username: $username, password: $password) {
      LibreLinkUpUsername
      LibreLinkUpPassword
      LibreLinkUpRegion
    }
  }
`);

function App() {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const { loading, error, data } = useQuery(GET_SETTINGS);
  const [saveSettings, { error: saveError }] = useMutation(SAVE_SETTINGS);

  useEffect(() => {
    if (!loading && data) {
      setUsername(data.settings.LibreLinkUpUsername);
      setPassword(data.settings.LibreLinkUpPassword);

    }
    // if (!saveLoading && saveData) {
    //   setUsername(saveData.settings.LibreLinkUpUsername);
    //   setPassword(saveData.settings.LibreLinkUpPassword);
    // }
  }, [loading, data]);

  // if (saveLoading) return <p>Loading...</p>;
  // if (saveError) return <p>Error : {saveError.message}</p>;

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error : {error.message}</p>;

  function onChange(username: string, password: string) {
    setUsername(username);
    setPassword(password);
  }

  function onSave() {
    saveSettings({
      variables: {
        username: username,
        password: password,
      }
    });
    // saveSettings({
    //   variables: {
    //     username: username,
    //     password: password,
    //   }
    // })
    //   .then((response) => {
    //     console.log("response:", response);

    //     const savedSettings = response.data?.saveSettings;
    //     if (savedSettings) {
    //       setUsername(savedSettings.LibreLinkUpUsername);
    //       setPassword(savedSettings.LibreLinkUpPassword);
    //       console.log("region: ", savedSettings.LibreLinkUpRegion);
    //     }
    //   })
    //   .catch((error) => {
    //     console.log("error:", error.message);
    //   });
  }

  return (
    <>
      {
        data &&
        <div className="">
          <LibreLinkupForm
            username={username}
            password={password}
            errorMessage={saveError?.message}
            onChange={onChange}
            onSave={onSave}
          />
        </div>
      }
    </>
  )
}

interface LibreLinkupFormProps {
  username: string
  password: string
  errorMessage: string | undefined
  onChange: (username: string, password: string) => void
  onSave: () => void
}

function LibreLinkupForm(props: LibreLinkupFormProps) {
  return (
    <div className="p-6 max-w-sm mx-auto bg-white rounded-xl shadow-lg items-center space-x-4">
      <div className="shrink-0 text-xl font-medium mb-4">
        LibreLinkUp Settings
      </div>
      <div className="bg-red-400">
        {props.errorMessage}
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
          onChange={(e) => props.onChange(e.target.value, props.password)}
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
          onChange={(e) => props.onChange(props.username, e.target.value)}
          placeholder="Password"
        />
      </div>
      <button type="button" onClick={props.onSave} >Save</button>
    </div>
  )
}

export default App
