package vgui

func FromFileRootPath(file string, rootPath string) (*Object, error) {
    obj, err := FromFileSourceProvider(file, &HudFileSourceProvider{rootPath})
    return obj, err
}

func FromFileSourceProvider(file string, provider FileSourceProvider) (*Object, error) {
    processed, err := Preprocess(provider, file)
    if err != nil {
        return nil, err
    }

    lexer, err := NewLexerFromInput(processed)
    if err != nil {
        return nil, err
    }

    parser, err := NewParserFromTokens(lexer.GetTokens())
    if err != nil {
        return nil, err
    }

    root, err := parser.ParseRoot()
    if err != nil {
        return nil, err
    }

    return FromValue(root), nil
}

func FromValue(value *Value) *Object {
    if value.String != nil {
        return NewObjectFromValue(value.Name.Value, value.String.Value)
    }
    obj := NewObject(value.Name.Value)
    for _, subValue := range value.Body {
        obj.mergeOrAddProperty(FromValue(subValue))
    }
    return obj
}