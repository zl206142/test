package util

type Model struct {
	TexturesLoaded  []Texture
	Meshes          []Mesh
	Directory       string
	GammaCorrection bool
}

func (model *Model) Draw(shader Shader) {
	for i := 0; i < len(model.Meshes); i++ {
		model.Meshes[i].Draw(shader)
	}
}
func (model *Model) loadModel(s string) {
	//todo:load
	/*
	// read file via ASSIMP
        Assimp::Importer importer;
        const aiScene* scene = importer.ReadFile(path, aiProcess_Triangulate | aiProcess_FlipUVs | aiProcess_CalcTangentSpace);
        // check for errors
        if(!scene || scene->mFlags & AI_SCENE_FLAGS_INCOMPLETE || !scene->mRootNode) // if is Not Zero
        {
            cout << "ERROR::ASSIMP:: " << importer.GetErrorString() << endl;
            return;
        }
        // retrieve the directory path of the filepath
        directory = path.substr(0, path.find_last_of('/'));

        // process ASSIMP's root node recursively
        processNode(scene->mRootNode, scene);

	*/
}

func NewModel(path string,gamma bool)*Model  {
	model:=Model{}
	model.GammaCorrection = gamma
	model.loadModel(path)
	return &model
}
