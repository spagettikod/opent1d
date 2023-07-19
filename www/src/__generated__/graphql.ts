/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type Query = {
  __typename?: 'Query';
  settings: Settings;
};

export type Settings = {
  __typename?: 'Settings';
  LibreLinkUpEndpoint: Scalars['String']['output'];
  LibreLinkUpPassword: Scalars['String']['output'];
  LibreLinkUpUsername: Scalars['String']['output'];
};

export type RootQueryQueryVariables = Exact<{ [key: string]: never; }>;


export type RootQueryQuery = { __typename?: 'Query', settings: { __typename?: 'Settings', LibreLinkUpUsername: string, LibreLinkUpEndpoint: string } };


export const RootQueryDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"RootQuery"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"settings"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"LibreLinkUpUsername"}},{"kind":"Field","name":{"kind":"Name","value":"LibreLinkUpEndpoint"}}]}}]}}]} as unknown as DocumentNode<RootQueryQuery, RootQueryQueryVariables>;