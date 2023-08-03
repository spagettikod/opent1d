import { gql, useMutation, useQuery } from "@apollo/client";
import { Button, Card, CardBody, CardFooter, CardHeader, FormControl, FormLabel, Heading, Input, Spacer, Stack, Text, useToast } from "@chakra-ui/react";
import { useEffect, useState } from "react";

const GET_SETTINGS = gql(`
  query GetSettings {
    settings {
      LibreLinkUpUsername
    }
  }
`);

const SAVE_SETTINGS = gql(`
  mutation SaveSettings($username: String!, $password: String!) {
    saveSettings(username: $username, password: $password) {
      LibreLinkUpUsername
    }
  }
`);

function LibreLinkUpDialog() {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const { loading, error, data } = useQuery(GET_SETTINGS);
    const [saveSettings] = useMutation(SAVE_SETTINGS);
    const toast = useToast();

    useEffect(() => {
        if (!loading && data) {
            setUsername(data.settings.LibreLinkUpUsername);
        }
    }, [loading, data]);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    function onSave() {
        saveSettings({
            variables: {
                username: username,
                password: password,
            }
        }).then((result) => {
            const response = result.data.saveSettings;
            if (response) {
                setUsername(response.LibreLinkUpUsername);
                toast({
                    title: "Sign in to LibreLinkUp",
                    description: "Login was successfull!",
                    status: 'success',
                    duration: 9000,
                    isClosable: true,
                    position: 'top'
                })
            }
        }).catch((error) => {
            console.error(error);
            toast({
                title: "Sign in to LibreLinkUp",
                description: error.message,
                status: 'error',
                duration: 9000,
                isClosable: true,
                position: 'top'
            })
        })
    }

    return (
        <>
            <Card alignContent="start">
                <CardHeader fontSize="xl" textAlign={"left"} fontWeight={"bold"}>
                    <Stack spacing={"2"}>
                        <Heading>LibreLinkUp</Heading>
                        <Text fontSize="lg" fontWeight={"thin"} textAlign={"left"}>
                            Enter your LibreLinkUp credentials to keep your glucose data in sync.
                        </Text>
                    </Stack>
                </CardHeader>
                <CardBody>
                    <Stack spacing={"3"}>
                        <FormControl>
                            <FormLabel>Username</FormLabel>
                            <Input type="text" value={username} onChange={(event) => setUsername(event.target.value)} />
                        </FormControl>
                        <FormControl>
                            <FormLabel>Password</FormLabel>
                            <Input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
                        </FormControl>
                    </Stack>
                </CardBody>
                <CardFooter>
                    <Spacer minWidth={'fit-content'} />
                    <Button
                        colorScheme="blue"
                        onClick={onSave}>
                        Save
                    </Button>
                </CardFooter>
            </Card>
        </>
    )

}

export default LibreLinkUpDialog