package objects

import (
	"bufio"
	"gonum.org/v1/gonum/mat"
	"os"
	"strconv"
	"strings"
)

func ReadFromObj(path string) VertexObject {

	object := VertexObject{
		Faces: make([]*Triangle, 0),
		Transformations: mat.NewDense(4, 4, []float64{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		}),
	}

	vertices := make([]*mat.Dense, 0)
	normals := make([]*mat.VecDense, 0)
	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := string([]rune(scanner.Text())[0:2])
		switch text {
		case "v ":
			rawSplit := strings.Split(string([]rune(scanner.Text())[2:]), " ")
			first, _ := strconv.ParseFloat(rawSplit[0], 64)
			second, _ := strconv.ParseFloat(rawSplit[1], 64)
			third, _ := strconv.ParseFloat(rawSplit[2], 64)
			vertices = append(vertices, mat.NewDense(4, 1, []float64{first, second, third, 1}))
		case "vn":
			rawSplit := strings.Split(string([]rune(scanner.Text())[2:]), " ")
			first, _ := strconv.ParseFloat(rawSplit[0], 64)
			second, _ := strconv.ParseFloat(rawSplit[1], 64)
			third, _ := strconv.ParseFloat(rawSplit[2], 64)
			normals = append(normals, mat.NewVecDense(3, []float64{first, second, third}))
		case "f ":
			initialSplit := strings.Split(strings.Replace(string([]rune(scanner.Text())[2:]), "//", "/-/", -1), " ")
			typeSplit0 := strings.Split(initialSplit[0], "/")
			typeSplit1 := strings.Split(initialSplit[1], "/")
			typeSplit2 := strings.Split(initialSplit[2], "/")
			firstV, _ := strconv.ParseInt(typeSplit0[0], 10, 64)
			secondV, _ := strconv.ParseInt(typeSplit1[0], 10, 64)
			thirdV, _ := strconv.ParseInt(typeSplit2[0], 10, 64)
			firstN, _ := strconv.ParseInt(typeSplit0[2], 10, 64)
			secondN, _ := strconv.ParseInt(typeSplit1[2], 10, 64)
			thirdN, _ := strconv.ParseInt(typeSplit2[2], 10, 64)
			object.Faces = append(object.Faces, &Triangle{
				RawVertices: [3]*mat.Dense{
					vertices[firstV-1],
					vertices[secondV-1],
					vertices[thirdV-1],
				},
				VertNormals: [3]*mat.VecDense{
					normals[firstN-1],
					normals[secondN-1],
					normals[thirdN-1],
				},
				FaceNormal: mat.NewVecDense(3, []float64{
					(normals[firstN-1].At(0, 0) + normals[secondN-1].At(0, 0) + normals[thirdN-1].At(0, 0)) / 3,
					(normals[firstN-1].At(1, 0) + normals[secondN-1].At(1, 0) + normals[thirdN-1].At(1, 0)) / 3,
					(normals[firstN-1].At(2, 0) + normals[secondN-1].At(2, 0) + normals[thirdN-1].At(2, 0)) / 3,
				}),
				TransformedVertices: [3]*mat.Dense{
					mat.NewDense(4, 1, nil),
					mat.NewDense(4, 1, nil),
					mat.NewDense(4, 1, nil),
				},
				Visibility: FRONT,
			})
		}
	}

	return object

}
